package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *TasksRepo) CreateTask(userID, description string) error {
	q := sq.
		Insert("tasks").
		Columns("task_id", "user_id", "description").
		Values(uuid.New().String(), userID, description)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *TasksRepo) UpdateTaskStatus(taskID, status int) error {
	q := sq.
		Update("tasks").
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"task_id": taskID})

	_, err := q.RunWith(r.db).Exec()
	return err
}

func (r *TasksRepo) UpdateTaskUser(taskID, userID string) error {
	q := sq.
		Update("tasks").
		Set("user_id", userID).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"task_id": taskID})

	_, err := q.RunWith(r.db).Exec()
	return err
}

func (r *TasksRepo) GetRandomUserID() (*uuid.UUID, error) {
	q := sq.
		Select("user_id").
		From("users").
		OrderBy("RAND()").
		Limit(1).
		Where(
			sq.Eq{"role_id": 4},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return &uuid.UUID{}, err
	}

	userID := &uuid.UUID{}
	err = r.db.Get(userID, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return &uuid.UUID{}, err
	}

	return userID, nil
}
