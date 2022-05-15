package controller

import (
	"encoding/json"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TasksController struct {
	tasksUsecase    *usecase.TasksUsecase
	rabbitmqService *rabbitmq.RabbitmqService
}

func NewTasksController(tasksUsecase *usecase.TasksUsecase, rabbitmqService *rabbitmq.RabbitmqService) *TasksController {
	return &TasksController{
		tasksUsecase:    tasksUsecase,
		rabbitmqService: rabbitmqService,
	}
}

// HTTP
func (c *TasksController) HandleCreateNewTask(context *gin.Context) {
	//c.tasksUsecase.CreateTask(c.)

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.TASK_CREATED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "tasks_service",
		Data: map[string]interface{}{
			"public_id": "",
			//todo: other info
		},
	}

	isValid := e.Validate(event.TASK_CREATED, 1)
	if isValid {
		c.rabbitmqService.Send(e, "new_tasks")
	} else {
		//TODO: retry or consume
	}

	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Added new task",
	})
}
func (c *TasksController) HandleShuffleTasks(context *gin.Context) {
	//c.tasksUsecase.ShuffleTasks()

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.TASK_CREATED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "tasks_service",
		Data: map[string]interface{}{
			//TODO: shuffled tasks array ?
		},
	}

	isValid := e.Validate(event.TASK_CREATED, 1)
	if isValid {
		c.rabbitmqService.Send(e, "shuffled_tasks")
	} else {
		//TODO: retry or consume
	}

	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "All open task are shuffled",
	})
}

//func (c *TasksController) HandleStatus(context *gin.Context) {
//	data, exists := context.Get("userdata")
//	rm := data.(*model.ResponseMessage)
//
//	if !exists {
//		context.JSON(http.StatusBadRequest, model.ResponseMessage{
//			Msg:  "Userdata not exists",
//			Code: http.StatusBadRequest,
//		})
//	}
//
//	context.JSON(http.StatusOK, model.ResponseMessage{
//		Msg:  fmt.Sprintf("PublicID: %s | Username: %s", rm.Claims.PublicID, rm.Claims.Username),
//		Code: http.StatusOK,
//	})
//}

// MQ
func (c *TasksController) HandleEvents(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.USER_CREATED:
		c.CreateUserWithBalance(e)
	case event.USER_UPDATED:
		log.Println("[Updated]", e.ID, e.Name, e.Version)
	case event.USER_DELETED:
		log.Println("[Deleted]", e.ID, e.Name, e.Version)
	}
}

func (c *TasksController) CreateUserWithBalance(e *rabbitmq.Event) {
	//TODO: create user with balance
	log.Println("[Created]", e.ID, e.Name, e.Version)
}
