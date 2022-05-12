package controller

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

	//TODO: remove event from here. Test
	e := rabbitmq.Event{
		EventID:      uuid.New().String(),
		EventVersion: 1,
		EventName:    "", //TODO: BEvent
		EventTime:    time.Now().Unix(),
		Producer:     "",
		Data:         nil,
	}

	//ac.rabbitmqService.ValidateEvent(e, "../../internal/event/schemas/tasks/created/1.json")
	ac.rabbitmqService.Send(e)

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
