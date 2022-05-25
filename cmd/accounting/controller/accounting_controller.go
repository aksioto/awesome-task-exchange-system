package controller

import (
	"encoding/json"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	model2 "github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
	v2 "github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v2"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type AccountingController struct {
	accountingUsecase *usecase.AccountingUsecase
	rabbitmqService   *rabbitmq.RabbitmqService
}

func NewAccountingController(accountingUsecase *usecase.AccountingUsecase, rabbitmqService *rabbitmq.RabbitmqService) *AccountingController {
	return &AccountingController{
		accountingUsecase: accountingUsecase,
		rabbitmqService:   rabbitmqService,
	}
}

func (c *AccountingController) HandleAccountingDashboard(ctx *gin.Context) {
	userData, exists := ctx.Get("userdata")
	rm := userData.(*model2.Claims)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to get dashboard.",
		})
		return
	}

	//TODO: move role check to middleware
	var data interface{}
	var err error

	switch rm.RoleID {
	case model2.ROLE_USER:
		data, err = c.accountingUsecase.GetUserBalanceWithLogForLast24Hours(rm.PublicID)
	case model2.ROLE_ADMIN, model2.ROLE_ACCOUNTANT:
		data, err = c.accountingUsecase.GetAccountingStatistic()
	default:
		err = errors.New("You shall not pass!")
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to get dashboard." + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Accounting dashboard",
		"data": data,
	})
}
func (c *AccountingController) HandleClosBillingCycle(ctx *gin.Context) {
	// send mail, update balance
	c.accountingUsecase.CloseBillingCycle()

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Task completed",
	})
}

// MQ
func (c *AccountingController) HandleTaskStream(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.TASK_CREATED:
		if e.Version == 1 {
			c.createTaskV1(e)
		} else if e.Version == 2 {
			c.createTaskV2(e)
		}
	case event.TASK_UPDATED:
		//c.updateTask(e)
	case event.TASK_DELETED:
		//c.deleteTask(e)
	}
}

func (c *AccountingController) createTaskV1(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.TASK_CREATED, 1); !isValid {
		log.Printf("%s", err)
		return
	}

	data := &v1.TaskData{}
	err := helper.MapToStruct(e.Data, data)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.CreateTaskV1(data)
	if err != nil {
		log.Println("Failed to create task. ", err.Error())
	} else {
		log.Println("[Created]", e.Name, e.Version, data.Title, data.PublicTaskID)
	}
}
func (c *AccountingController) createTaskV2(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.TASK_CREATED, 2); !isValid {
		log.Printf("%s", err)
		return
	}

	data := &v2.TaskData{}
	err := helper.MapToStruct(e.Data, data)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.CreateTaskV2(data)
	if err != nil {
		log.Println("Failed to create task. ", err.Error())
	} else {
		log.Println("[Created]", e.Name, e.Version, data.Title, data.JiraID, data.PublicTaskID)
	}
}

func (c *AccountingController) HandleUserStream(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.USER_CREATED:
		c.createUser(e)
	case event.USER_UPDATED:
		//c.updateUser(e)
	case event.USER_DELETED:
		//c.deleteUser(e)
	}
}

func (c *AccountingController) createUser(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.USER_CREATED, 1); !isValid {
		log.Printf("%s", err)
		return
	}

	userData := &v1.UserData{}
	err := helper.MapToStruct(e.Data, userData)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.CreateUserV1(userData)
	if err != nil {
		log.Println("Failed to create user. ", err.Error())
	} else {
		log.Println("[Created]", e.Name, e.Version, userData.Name, userData.PublicID)
	}
}

func (c *AccountingController) HandleTaskStatuses(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.TASK_COMPLETED:
		c.completeTask(e)
	}
}
func (c *AccountingController) completeTask(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.TASK_COMPLETED, 1); !isValid {
		log.Printf("%s", err)
		return
	}

	data := &v1.TaskAssigneeData{}
	err := helper.MapToStruct(e.Data, data)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.CompleteTask(data)
	if err != nil {
		log.Println("Failed to complete. ", err.Error())
	} else {
		log.Println("[Completed]", e.Name, e.Version)
	}

	_ = c.accountingUsecase.ApplyTransaction(data.PublicUserID, data.PublicTaskID, "increase", "")
}

func (c *AccountingController) HandleTaskAssignment(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.TASKS_SHUFFLED:
		c.shuffleTask(e)
	case event.TASK_ASSIGNED:
		c.assignTask(e)
	}
}
func (c *AccountingController) shuffleTask(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.TASKS_SHUFFLED, 1); !isValid {
		log.Printf("%s", err)
		return
	}

	data := &v1.TaskAssigneeData{}
	err := helper.MapToStruct(e.Data, data)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.AssignTaskV1(data)
	if err != nil {
		log.Println("Failed to shuffle. ", err.Error())
	} else {
		log.Println("[Shuffled]", e.Name, e.Version)
	}

	_ = c.accountingUsecase.ApplyTransaction(data.PublicUserID, data.PublicTaskID, "decrease", "")
}

func (c *AccountingController) assignTask(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.TASK_ASSIGNED, 1); !isValid {
		log.Printf("%s", err)
		return
	}
	data := &v1.TaskAssigneeData{}
	err := helper.MapToStruct(e.Data, data)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.accountingUsecase.AssignTaskV1(data)
	if err != nil {
		log.Println("Failed to assign. ", err.Error())
	} else {
		log.Println("[Assigned]", e.Name, e.Version)
	}

	_ = c.accountingUsecase.ApplyTransaction(data.PublicUserID, data.PublicTaskID, "decrease", "")
}
