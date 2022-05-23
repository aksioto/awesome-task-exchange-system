package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

type TasksRepo struct {
	db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
	return &TasksRepo{
		db: db,
	}
}

func (r *TasksRepo) CreateTask(title, description string) (*model.Task, error) {
	q := sq.
		Insert("tasks").
		Columns("task_id", "title", "description").
		Values(uuid.New().String(), title, description)

	res, err := q.RunWith(r.db).Exec()
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetTaskByInternalID(lastInsertID)
}

func (r *TasksRepo) GetTaskByInternalID(internalID int64) (*model.Task, error) {
	q := sq.
		Select("*").
		From("tasks").
		Where(
			sq.Eq{"id": internalID},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	task := &model.Task{}
	err = r.db.Get(task, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return nil, err
	}

	return task, nil
}

func (r *TasksRepo) UpdateTaskStatus(taskID string, status int) error {
	q := sq.
		Update("tasks").
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(
			sq.And{
				sq.Eq{"task_id": taskID},
				sq.Eq{"status": 0},
			})

	res, err := q.RunWith(r.db).Exec()
	rows, err := res.RowsAffected()
	if rows == 0 {
		return errors.New("Nothing to update")
	}
	return err
}

func (r *TasksRepo) UpdateTaskUser(taskID, userID uuid.UUID) error {
	q := sq.
		Update("tasks").
		Set("user_id", userID.String()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"task_id": taskID.String()})

	_, err := q.RunWith(r.db).Exec()
	return err
}

func (r *TasksRepo) GetTasksIDsWithZeroStatus() ([]uuid.UUID, error) {
	q := sq.
		Select("task_id").
		From("tasks").
		Where(
			sq.Eq{"status": 0},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	var tasksIDs []uuid.UUID
	err = r.db.Select(&tasksIDs, sqlQ, args...)
	if err != nil {
		log.Printf("DB1: %s", err.Error())
		return nil, err
	}

	return tasksIDs, nil
}

func (r *TasksRepo) GetRandomUserID() (uuid.UUID, error) {
	q := sq.
		Select("public_id").
		From("users").
		OrderBy("RAND()").
		Limit(1).
		Where(
			sq.Eq{"role_id": 4},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return uuid.UUID{}, err
	}

	userID := uuid.UUID{}
	err = r.db.Get(&userID, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (r *TasksRepo) AddRandomUserToTask(taskID uuid.UUID) (uuid.UUID, error) {
	randomUserID, err := r.GetRandomUserID()
	if err != nil {
		return uuid.UUID{}, err
	}

	return randomUserID, r.UpdateTaskUser(taskID, randomUserID)
}

func (r *TasksRepo) CreateUserV1(userData *v1.UserData) error {
	q := sq.
		Insert("users").
		Columns("public_id", "email", "name", "role_id").
		Values(userData.PublicID.String(), userData.Email, userData.Name, userData.RoleID)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *TasksRepo) GetAssignedTasks(userID uuid.UUID) ([]*model.Task, error) {
	q := sq.
		Select("*").
		From("tasks").
		Where(
			sq.And{
				sq.Eq{"user_id": userID.String()},
				sq.Eq{"status": 0},
			},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	var tasks []*model.Task
	err = r.db.Select(&tasks, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return nil, err
	}

	return tasks, nil
}
