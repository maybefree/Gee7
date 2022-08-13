package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	// 1.创建路由
	r := gin.Default()
	r.Use(Logger())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})

	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello World!")
	})

	r.Run(":8000")
}