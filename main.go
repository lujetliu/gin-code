package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.RecoveryFunc {
	return gin.RecoveryFunc(func(c *gin.Context, err interface{}) {
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		headers := strings.Split(string(httpRequest), "\r\n")
		for idx, header := range headers {
			current := strings.Split(header, ":")
			if current[0] == "Authorization" {
				headers[idx] = current[0] + ": *"
			}
		}
		headersToStr := strings.Join(headers, "\r\n")
		fmt.Println(headersToStr)
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func main() {
	r := gin.Default()
	r.Use(gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, Recovery()))

	r.GET("/ping", func(c *gin.Context) {
		v := []int{1, 2, 3}
		_ = v[5]
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	r.Run("127.0.0.1:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
