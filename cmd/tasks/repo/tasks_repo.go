package repo

import "github.com/jmoiron/sqlx"

type TasksRepo struct {
	db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
	return &TasksRepo{
		db: db,
	}
}
