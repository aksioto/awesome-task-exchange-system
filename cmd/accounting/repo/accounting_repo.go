package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountingRepo struct {
	db *sqlx.DB
}

func NewAccountingRepo(db *sqlx.DB) *AccountingRepo {
	return &AccountingRepo{
		db: db,
	}
}

func (r *AccountingRepo) AddLog(userID, description string) error {
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
