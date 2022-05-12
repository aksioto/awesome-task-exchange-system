package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/repo"
)

type TasksUsecase struct {
	tasksRepo *repo.TasksRepo
}

func NewTasksUsecase(tasksRepo *repo.TasksRepo) *TasksUsecase {
	return &TasksUsecase{
		tasksRepo: tasksRepo,
	}
}

func (tu *TasksUsecase) CreateTask() {
	//tu.tasksRepo.CreateTask()
}
