package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"realTimeService/configuration"
	"realTimeService/handlers"
	"realTimeService/interfaces"
	"realTimeService/middlewares"
	"realTimeService/providers"
)

func main() {
	// Load configuration
	cfg, err := configuration.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	router := gin.Default()
	// Initialize the DI container
	var container interfaces.Container = providers.NewDependencyInjectionContainer()
	container.InitializeProviders(cfg)
	defer func(container interfaces.Container) {
		err := container.Close()
		if err != nil {
			log.Fatalf("Failed to close dependency injection container: %v", err)
			return
		}
	}(container)
	wsHandler := handlers.NewWsHandler(container)
	router.Use(gin.Recovery())
	router.GET("/ws", middlewares.AuthMiddleware(cfg), wsHandler.Handle)
	err = router.Run(cfg.HttpPort)
	if err != nil {
		return
	}
}
