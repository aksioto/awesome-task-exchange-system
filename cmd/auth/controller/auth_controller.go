package controller

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHeader struct {
	Authorization string `header:"Authorization"`
}

type AuthController struct {
	authUsecase     *usecase.AuthUsecase
	rabbitmqService *rabbitmq.RabbitmqService
}

func NewAuthController(authUsecase *usecase.AuthUsecase, rabbitmqService *rabbitmq.RabbitmqService) *AuthController {
	return &AuthController{
		authUsecase:     authUsecase,
		rabbitmqService: rabbitmqService,
	}
}

func (ac *AuthController) HandleSignIn(c *gin.Context) {
	email := c.PostForm("email")
	pass := c.PostForm("password")

	token, _, err := ac.authUsecase.SignIn(email, pass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to sign in. Check your email or password.",
		})
		return
	}

	c.SetCookie("token", token, int(time.Hour), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Successfully signed in",
	})
}

func (ac *AuthController) HandleToken(c *gin.Context) {
	header := &authHeader{}
	if err := c.ShouldBindHeader(header); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	token := strings.Split(header.Authorization, "token ")
	claims, err := ac.authUsecase.VerifyToken(token[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"msg":    "Token is valid",
		"claims": &claims,
	})
}

func (ac *AuthController) HandleSignup(c *gin.Context) {
	//TODO: signup logic

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.USER_CREATED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "auth_service",
		Data: map[string]interface{}{
			"public_id": "",
			//todo: other user info
		},
	}

	isValid := e.Validate(event.USER_CREATED, 1)
	if isValid {
		ac.rabbitmqService.Send(e, "")
	} else {
		//TODO: retry or send error log
	}
}
