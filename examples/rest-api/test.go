package main

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" :"pong",
		})
	})

	r.Run()
}
