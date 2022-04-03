package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	connectMySQL()
	RedisNewClient()
	server := gin.Default()

	rand.Seed(time.Now().UnixNano())

	server.POST("/api/v1/urls", rest)
	server.GET("/:id", redirect)
	server.GET("/", index)
	server.Run(":80")
}
