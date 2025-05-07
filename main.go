package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"realTimeService/clients"
	"realTimeService/configuration"
	"realTimeService/handlers"
	"realTimeService/hubs"
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
	grpcClient := clients.NewMessageServiceClient(conn)
	if grpcClient == nil {
		log.Fatalf("gRPC client creation failed")
		return
	}
	hub := hubs.NewMainHub()
	router.GET("/ws", func(c *gin.Context) {
		handlers.WsHandler(c, hub, grpcClient)
	})
	err = router.Run(cfg.HttpPort)
	if err != nil {
		return
	}
}
