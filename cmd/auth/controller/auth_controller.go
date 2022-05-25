package controller

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
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

func (c *AuthController) HandleSignIn(ctx *gin.Context) {
	email := ctx.PostForm("email")
	pass := ctx.PostForm("password")

	token, err := c.authUsecase.SignIn(email, pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to sign in. Check your email or password.",
		})
		return
	}

	ctx.SetCookie("token", token, int(time.Hour), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Successfully signed in",
	})
}

func (c *AuthController) HandleToken(ctx *gin.Context) {
	header := &authHeader{}
	if err := ctx.ShouldBindHeader(header); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	token := strings.Split(header.Authorization, "token ")
	claims, err := c.authUsecase.VerifyToken(token[1])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"msg":    "Token is valid",
		"claims": &claims,
	})
}

func (c *AuthController) HandleSignUp(ctx *gin.Context) {
	user, err := c.authUsecase.SignUp(ctx.PostForm("email"), ctx.PostForm("password"), ctx.PostForm("name"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	e := rabbitmq.Event{
		ID:       uuid.New().String(),
		Version:  1,
		Name:     event.USER_CREATED,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		Producer: "auth_service",
		Data: map[string]interface{}{
			"public_id": user.PublicID,
			"email":     user.Email,
			"name":      user.Name,
			"role_id":   user.RoleID,
		},
	}

	isValid, err := e.Validate(event.USER_CREATED, 1)
	if isValid {
		log.Printf("Valid event. %s", user.Email)
		_ = c.rabbitmqService.Send(e.ToJson(), "user_stream")

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "User created",
		})
	} else {
		//TODO: retry or send error log

		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Event validation failed. " + err.Error(),
		})
	}
}
