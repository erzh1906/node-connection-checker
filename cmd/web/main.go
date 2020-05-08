package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d %s\n",
			params.ClientIP,
			params.Request.Header["Request-Id"],
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
			)
	}))
	router.Use(gin.Recovery())
	router.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	err := router.Run(":20000")
	if err != nil {
		log.Fatal(err)
	}
}
