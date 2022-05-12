package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/repo"
	"github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"time"
)

type AuthUsecase struct {
	authRepo     *repo.AuthRepo
	clientSecret []byte
}

func NewAuthUsecase(clientSecret string, authRepo *repo.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo:     authRepo,
		clientSecret: []byte(clientSecret),
	}
}

func (u *AuthUsecase) SignIn(email, pass string) (string, string, error) {
	user, err := u.authRepo.GetUser(email, pass)
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	token, err := u.generateToken(user)
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	if err = u.authRepo.SaveAuthToken(user.PublicID.String(), token); err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	return token, user.PublicID.String(), nil
}

func (u *AuthUsecase) generateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := &model.Claims{
		Username: user.Name,
		PublicID: user.PublicID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.clientSecret)
}

func (u *AuthUsecase) VerifyToken(token string) (*model.Claims, error) {
	claims, err := u.parseToken(token)
	if err != nil {
		return nil, err
	}

	count, err := u.authRepo.GetUserToken(token, claims.PublicID.String())
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

func (u *AuthUsecase) parseToken(token string) (*model.Claims, error) {
	claims := &model.Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return u.clientSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
