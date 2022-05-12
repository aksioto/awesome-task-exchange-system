package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/repo"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
)

type TasksUsecase struct {
	tasksRepo       *repo.TasksRepo
	rabbitmqService *rabbitmq.RabbitmqService
}

func NewTasksUsecase(tasksRepo *repo.TasksRepo, rabbitmqService *rabbitmq.RabbitmqService) *TasksUsecase {
	return &TasksUsecase{
		tasksRepo:       tasksRepo,
		rabbitmqService: rabbitmqService,
	}
}
