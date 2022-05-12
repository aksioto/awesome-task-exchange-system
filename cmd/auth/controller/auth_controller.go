package controller

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type authHeader struct {
	Authorization string `header:"Authorization"`
}

type AuthController struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthController(authUsecase *usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (uc *AuthController) HandleSignIn(c *gin.Context) {
	email := c.PostForm("email")
	pass := c.PostForm("password")

	token, err := uc.authUsecase.SignIn(email, pass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Failed to sign in. Check your email or password.",
		})
		return
	}

	//expirationTime := time.Now().Add(60 * time.Minute)
	c.SetCookie("token", token, int(60*time.Minute), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Successfully signed in",
	})
}

func (uc *AuthController) HandleToken(c *gin.Context) {
	header := &authHeader{}
	if err := c.ShouldBindHeader(header); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	token := strings.Split(header.Authorization, "token ")
	claims, err := uc.authUsecase.VerifyToken(token[1])
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
