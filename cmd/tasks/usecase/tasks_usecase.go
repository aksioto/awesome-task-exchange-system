package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/repo"
	"github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
	"github.com/google/uuid"
	"log"
)

type TasksUsecase struct {
	tasksRepo *repo.TasksRepo
}

func NewTasksUsecase(tasksRepo *repo.TasksRepo) *TasksUsecase {
	return &TasksUsecase{
		tasksRepo: tasksRepo,
	}
}

func (u *TasksUsecase) CreateUserV1(userData *v1.UserData) error {
	err := u.tasksRepo.CreateUserV1(userData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (u *TasksUsecase) CreateTask(title, description string) (*model.Task, error) {
	task, err := u.tasksRepo.CreateTask(title, description)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return task, nil
}

func (u *TasksUsecase) AssignTask(taskID uuid.UUID) (uuid.UUID, error) {
	userID, err := u.tasksRepo.AddRandomUserToTask(taskID)
	if err != nil {
		log.Println(err.Error())
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (u *TasksUsecase) ReshuffleTasks() ([]map[string]interface{}, error) {
	reshuffledTasks := make([]map[string]interface{}, 1)

	taskIDs, err := u.tasksRepo.GetTasksIDsWithZeroStatus()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for _, taskID := range taskIDs {
		rndUserID, err := u.tasksRepo.GetRandomUserID()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = u.tasksRepo.UpdateTaskUser(taskID, rndUserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		rt := make(map[string]interface{}, 2)
		rt["public_task_id"] = taskID.String()
		rt["public_user_id"] = rndUserID.String()
		reshuffledTasks = append(reshuffledTasks, rt)
	}

	return reshuffledTasks, nil
}

func (u *TasksUsecase) CompleteTask(taskID string) error {
	err := u.tasksRepo.UpdateTaskStatus(taskID, 1)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (u *TasksUsecase) GetAssignedTasks(userID uuid.UUID) ([]*model.Task, error) {
	tasks, err := u.tasksRepo.GetAssignedTasks(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return tasks, nil
}
