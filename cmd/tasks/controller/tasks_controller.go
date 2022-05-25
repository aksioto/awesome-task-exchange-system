package controller

import (
	"encoding/json"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	message "github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
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
func (c *TasksController) HandleAddNewTask(ctx *gin.Context) {
	task, err := c.tasksUsecase.CreateTask(ctx.PostForm("title"), ctx.PostForm("description"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to add new task. " + err.Error(),
		})
		return
	}

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.TASK_CREATED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "tasks_service",
		Data: map[string]interface{}{
			"public_task_id": task.TaskID,
			"title":          task.Title,
			"description":    task.Description,
		},
	}

	isValid, err := e.Validate(event.TASK_CREATED, 1)
	if isValid {
		_ = c.rabbitmqService.Send(e.ToJson(), "task_stream")
	} else {
		//TODO: retry or send error log
		log.Println("[ATTENTION] event validation failed. " + err.Error())
	}

	err = c.assignTaskToRandomUser(task)
	if err != nil {
		log.Println(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "New task added",
	})
}

func (c *TasksController) assignTaskToRandomUser(task *model.Task) error {
	userID, err := c.tasksUsecase.AssignTask(task.TaskID)
	if err != nil {
		return err
	}

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.TASK_ASSIGNED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "tasks_service",
		Data: map[string]interface{}{
			"public_task_id": task.TaskID,
			"public_user_id": userID,
		},
	}

	isValid, err := e.Validate(event.TASK_ASSIGNED, 1)
	if isValid {
		_ = c.rabbitmqService.Send(e.ToJson(), "task_assignment")
	} else {
		//TODO: retry or send error log
		log.Println("[ATTENTION] event validation failed. " + err.Error())
		return err
	}

	return nil
}

func (c *TasksController) HandleShuffleTasks(ctx *gin.Context) {
	reshuffledTasks, _ := c.tasksUsecase.ReshuffleTasks()

	for _, reshuffledTask := range reshuffledTasks {
		e := rabbitmq.Event{
			ID:       uuid.New().String(),
			Version:  1,
			Name:     event.TASKS_SHUFFLED,
			Time:     strconv.FormatInt(time.Now().Unix(), 10),
			Producer: "tasks_service",
			Data:     reshuffledTask,
		}

		isValid, err := e.Validate(event.TASKS_SHUFFLED, 1)
		if isValid {
			_ = c.rabbitmqService.Send(e.ToJson(), "task_assignment")
		} else {
			//TODO: retry or consume
			log.Println(err.Error())
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "All open tasks are shuffled",
	})
}

func (c *TasksController) HandleCompleteTask(ctx *gin.Context) {
	taskID := ctx.PostForm("task_id")

	data, exists := ctx.Get("userdata")
	rm := data.(*message.Claims)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to complete task.",
		})
	}

	err := c.tasksUsecase.CompleteTask(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to complete task. " + err.Error(),
		})
		return
	}

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.TASK_COMPLETED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "tasks_service",
		Data: map[string]interface{}{
			"public_task_id": taskID,
			"public_user_id": rm.PublicID,
		},
	}

	isValid, err := e.Validate(event.TASK_COMPLETED, 1)
	if isValid {
		_ = c.rabbitmqService.Send(e.ToJson(), "task_statuses")
	} else {
		//TODO: retry or consume
		log.Println(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Task completed",
	})
}

func (c *TasksController) HandleDashboard(ctx *gin.Context) {
	data, exists := ctx.Get("userdata")
	rm := data.(*message.Claims)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to get tasks.",
		})
		return
	}

	tasks, err := c.tasksUsecase.GetAssignedTasks(rm.PublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to get tasks.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Assigned tasks",
		"data": tasks,
	})
}

// --> CUD users
func (c *TasksController) HandleUserStream(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.USER_CREATED:
		if e.Version == 1 {
			c.createUserV1(e)
		}
	case event.USER_UPDATED:
		log.Println("[Updated]", e.Name, e.Version)
	case event.USER_DELETED:
		log.Println("[Deleted]", e.Name, e.Version)
	}
}

func (c *TasksController) createUserV1(e *rabbitmq.Event) {
	if isValid, err := e.Validate(event.USER_CREATED, 1); !isValid {
		log.Printf("%s", err)
		return
	}

	userData := &v1.UserData{}
	err := helper.MapToStruct(e.Data, userData)
	if err != nil {
		log.Println("Failed convert map to struct. ", err.Error())
	}

	err = c.tasksUsecase.CreateUserV1(userData)
	if err != nil {
		log.Println("Failed to create user. ", err.Error())
	} else {
		log.Println("[Created]", e.Name, e.Version, userData.Name, userData.PublicID)
	}
}

// <-- CUD users
