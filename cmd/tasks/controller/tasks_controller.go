package controller

import (
	"encoding/json"
	"fmt"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
func (tc *TasksController) HandleCreateNewTask(c *gin.Context) {
	//tc.tasksUsecase.CreateTask(c.)
	//TODO: dispatch e = task.added
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Added new task",
	})
}
func (tc *TasksController) HandleShuffleTasks(c *gin.Context) {

	//tc.tasksUsecase.
	//TODO: dispatch e = task.shuffled

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "All open tasks are shuffled",
	})
}

func (tc *TasksController) HandleStatus(c *gin.Context) {
	data, exists := c.Get("userdata")
	rm := data.(*model.ResponseMessage)

	if !exists {
		c.JSON(http.StatusBadRequest, model.ResponseMessage{
			Msg:  "Userdata not exists",
			Code: http.StatusBadRequest,
		})
	}

	c.JSON(http.StatusOK, model.ResponseMessage{
		Msg:  fmt.Sprintf("PublicID: %s | Username: %s", rm.Claims.PublicID, rm.Claims.Username),
		Code: http.StatusOK,
	})
}

// MQ
func (tc *TasksController) StartReceiver() {
	log.Println(" [*] Receiver started")
	tc.rabbitmqService.Receive(tc.receiveWorker)
}

func (tc *TasksController) receiveWorker(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}
	log.Println(e.EventName)
	//
	//switch e.EventName {
	//case event.USER_SIGNEDIN:
	//	log.Printf("User signed in! %s", message.Data)
	//	break
	//case event.USER_GOT_TOKEN:
	//	log.Printf("User got token! %s", message.Data)
	//	break
	//default:
	//	break
	//}
}
