package middleware

import (
	"fmt"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectUrl := fmt.Sprintf("%s%s", c.Request.Host, c.Request.RequestURI)
		authUrl := fmt.Sprintf("http://localhost:8081/app?redirectUrl=%s", redirectUrl)
		tokenUrl := "http://127.0.0.1:8081/token"

		cookie, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.Redirect(http.StatusFound, authUrl)
				c.Abort()
				//c.JSON(http.StatusBadRequest, gin.H{
				//	"msg": "Auth. Missing cookie with token.",
				//})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Auth. Bad request.",
			})
			return
		}

		header := http.Header{
			"Authorization": []string{fmt.Sprintf("token %s", cookie)},
		}
		claims, err := makeRequest(tokenUrl, header)
		if err != nil {
			log.Printf("Redirecting to %s", authUrl)
			c.Redirect(http.StatusFound, authUrl)
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
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

	if res.StatusCode == http.StatusUnauthorized {
		return "", errors.New(fmt.Sprintf("Error! Status code %o", res.StatusCode))
	}
	return string(body), nil
}
