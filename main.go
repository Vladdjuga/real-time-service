package main

import (
	"github.com/gin-gonic/gin"
	"realTimeService/handlers"
)

func main() {
	router := gin.Default()
	router.GET("/ws", handlers.WsHandler)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
