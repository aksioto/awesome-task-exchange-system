package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) SaveAuthToken(userID, token string) error {
	q := sq.
		Insert("tokens").
		Columns("user_id", "token").
		Values(userID, token)

	_, err := q.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetUserToken(token, publicID string) (int, error) {
	q := sq.
		Select("COUNT(token)").
		From("tokens").
		Where(
			sq.And{
				sq.Eq{"token": token},
				sq.Eq{"user_id": publicID},
			},
		)

	sqlQ, args, err := q.ToSql()
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return 0, err
	}

	var count int
	err = r.db.Get(&count, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return 0, err
	}

	return count, nil
}

func (r *AuthRepo) GetUser(email, pass string) (*model.User, error) {
	q := sq.
		Select("*").
		From("users").
		Where(
			sq.And{
				sq.Eq{"email": email},
				sq.Eq{"password": pass},
			},
		)

	sqlQ, args, err := q.ToSql()
	//sqlboilerplaite here https://github.com/volatiletech/sqlboiler
	if err != nil {
		log.Printf("Can't sql from query: %s", spew.Sdump(q))
		return nil, err
	}

	user := &model.User{}
	err = r.db.Get(user, sqlQ, args...)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return nil, err
	}

	return user, nil
}
