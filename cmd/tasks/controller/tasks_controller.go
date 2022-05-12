package controller

import (
	"encoding/json"
	"fmt"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TasksController struct {
	tasksUsecase *usecase.TasksUsecase
}

func NewTasksController(tasksUsecase *usecase.TasksUsecase) *TasksController {
	return &TasksController{
		tasksUsecase: tasksUsecase,
	}
}

func (tc *TasksController) HandleAddNewTask(c *gin.Context) {

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
	rm := &model.ResponseMessage{}
	data, exists := c.Get("claims")

	if exists {
		err := json.Unmarshal([]byte(fmt.Sprint(data)), &rm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "parsing claims failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "Status OK",
		"public_id": rm.Claims.PublicID,
		"name":      rm.Claims.Username,
	})
}
