package repo

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/internal/model"
	v1 "github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
	v2 "github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type AccountingRepo struct {
	db *sqlx.DB
}

func NewAccountingRepo(db *sqlx.DB) *AccountingRepo {
	return &AccountingRepo{
		db: db,
	}
}

func (r *AccountingRepo) CreateUserV1(userData *v1.UserData) error {
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

func (r *AccountingRepo) CreateTaskV1(data *v1.TaskData) error {
	q := sq.
		Insert("tasks").
		Columns("task_id", "title", "description").
		Values(data.PublicTaskID.String(), data.Title, data.Description)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountingRepo) CreateTaskV2(data *v2.TaskData) error {
	q := sq.
		Insert("tasks").
		Columns("task_id", "title", "jira_id", "description").
		Values(data.PublicTaskID.String(), data.Title, data.JiraID, data.Description)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}
func (r *AccountingRepo) UpdateTaskStatus(taskID uuid.UUID, status int) error {
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

func (r *AccountingRepo) GetTaskDescription(taskID uuid.UUID) (string, error) {
	q := sq.
		Select("description").
		From("tasks").
		Where(
			sq.Eq{"task_id": taskID.String()},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return "", err
	}

	var description string
	err = r.db.Get(&description, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return "", err
	}

	return description, nil
}

func (r *AccountingRepo) UpdateTaskUser(taskID, userID uuid.UUID) error {
	q := sq.
		Update("tasks").
		Set("user_id", userID.String()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"task_id": taskID.String()})

	_, err := q.RunWith(r.db).Exec()
	return err
}

func (r *AccountingRepo) IncreaseTransaction(userID, taskID uuid.UUID, description string, val, balance int) error {
	q := sq.
		Insert("transactions").
		Columns("user_id", "task_id", "description", "increase", "balance").
		Values(userID.String(), taskID.String(), description, val, balance)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}
func (r *AccountingRepo) DecreaseTransaction(userID, taskID uuid.UUID, description string, val, balance int) error {
	q := sq.
		Insert("transactions").
		Columns("user_id", "task_id", "description", "decrease", "balance").
		Values(userID.String(), taskID.String(), description, val, balance)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}
func (r *AccountingRepo) PayoutTransaction(userID uuid.UUID, description string, val int) error {
	q := sq.
		Insert("transactions").
		Columns("user_id", "description", "decrease", "balance").
		Values(userID.String(), description, val, 0)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountingRepo) GetUserBalance(userID uuid.UUID) (int, error) {
	q := sq.
		Select("balance").
		From("transactions").
		Where(
			sq.Eq{"user_id": userID.String()},
		).
		OrderBy("created_at DESC").
		Limit(1)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return 0, err
	}

	var balance int
	err = r.db.Get(&balance, sqlQ, args...)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		log.Printf("DB: %s", err.Error())
		return 0, err
	}

	return balance, nil
}
func (r *AccountingRepo) GetUserSumBalance(userID uuid.UUID) (int, error) {
	q := sq.
		Select("SUM((increase - decrease)) as balance").
		From("transactions").
		Where(
			sq.And{
				sq.Eq{"user_id": userID.String()},
				sq.Expr("created_at >= now() - INTERVAL 1 DAY"),
			},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return 0, err
	}

	var balance int
	err = r.db.Get(&balance, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return 0, err
	}

	return balance, nil
}

func (r *AccountingRepo) GetUserTransactions(userID uuid.UUID) ([]model.Transaction, error) {
	q := sq.
		Select("*").
		From("transactions").
		Where(
			sq.And{
				sq.Eq{"user_id": userID.String()},
				sq.Expr("created_at >= now() - INTERVAL 1 DAY"),
			},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	var transactions []model.Transaction
	err = r.db.Select(&transactions, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return nil, err
	}

	return transactions, nil
}

func (r *AccountingRepo) GetProfit() (int, error) {
	q := sq.
		Select("(SELECT SUM(cost) FROM tasks WHERE status = 1 AND updated_at >= now() - INTERVAL 1 DAY) + sum(decrease-increase) as profit").
		From("transactions").
		Where(
			sq.And{
				sq.Expr("created_at >= now() - INTERVAL 1 DAY"),
			},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return 0, err
	}

	var profit int
	err = r.db.Get(&profit, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return 0, err
	}

	return profit, nil
}

func (r *AccountingRepo) GetUsersIDs() ([]uuid.UUID, error) {
	q := sq.
		Select("public_id").
		From("users").
		Where(
			sq.Eq{"role_id": 4},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	var usersIDs []uuid.UUID
	err = r.db.Select(&usersIDs, sqlQ, args...)
	if err != nil {
		log.Printf("DB1: %s", err.Error())
		return nil, err
	}

	return usersIDs, nil
}
