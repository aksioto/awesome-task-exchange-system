package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	"github.com/aksioto/awesome-task-exchange-system/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redirectUrl := fmt.Sprintf("%s%s", ctx.Request.Host, ctx.Request.RequestURI)
		authUrl := fmt.Sprintf("http://localhost:8081/app?redirectUrl=%s", redirectUrl)
		tokenUrl := "http://127.0.0.1:8081/token"

		cookie, err := ctx.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println(err.Error())
				redirect(ctx, authUrl)
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Auth. Bad request." + err.Error(),
			})
			return
		}

		header := http.Header{
			"Authorization": []string{fmt.Sprintf("token %s", cookie)},
		}
		claims, err := makeRequest(tokenUrl, header)
		if err != nil {
			log.Println("token validation failed")
			redirect(ctx, authUrl)
			return
		}

		rm := &model.ResponseMessage{}
		err = json.Unmarshal([]byte(claims), &rm)
		if err != nil {
			log.Println("parsing claims failed")
			return
		}

		ctx.Set("userdata", rm.Claims)
		ctx.Next()
	}
}

func redirect(ctx *gin.Context, url string) {
	log.Printf("Redirecting to %s", url)
	ctx.Redirect(http.StatusFound, url)
	ctx.Abort()
}

func makeRequest(url string, header http.Header) (string, error) {
	res, err := helper.Get(url, header)
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Error! Status code %o", res.StatusCode))
	}

	return string(body), nil
}
