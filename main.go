package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"realTimeService/clients"
	"realTimeService/configuration"
	"realTimeService/handlers"
	"realTimeService/hubs"
	"realTimeService/middlewares"
)

func main() {
	// Load configuration
	cfg, err := configuration.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	router := gin.Default()
	conn, err := grpc.Dial(cfg.GrpcClientAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("gRPC connect failed: %v", err)
		return
	}
	defer conn.Close()
	grpcMessagesClient := clients.NewMessageServiceClient(conn)
	if grpcMessagesClient == nil {
		log.Fatalf("gRPC client creation failed")
		return
	}
	hub := hubs.NewMainHub()
	router.Use(gin.Recovery())
	router.GET("/ws", middlewares.AuthMiddleware(cfg), func(c *gin.Context) {
		handlers.WsHandler(c, hub, grpcMessagesClient)
	})
	err = router.Run(cfg.HttpPort)
	if err != nil {
		return
	}
}
