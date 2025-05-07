package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"realTimeService/clients"
	"realTimeService/hubs"
	"realTimeService/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(c *gin.Context, hub *hubs.MainHub, grpcClient *clients.MessageServiceClient) {
	log.Println("WsHandler called.")
	// This is a placeholder for authentication.
	// In a real application, you would check the user's authentication token here.
	// For this example, we assume the user is authenticated and has a userId.
	userId := c.Query("userId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	var client models.Client
	if err != nil {
		log.Println("Write error:", err)
		return
	} else {
		log.Println("Connected to websocket")
		client = *models.NewClient(uuid.MustParse(userId), nil, conn)
	}
	defer conn.Close()

	for {
		// Read message from the client
		log.Println("Waiting for message...")
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Println("Message received:", string(msgBytes))
		// Handle the message
		HandleMessage(hub, grpcClient, client, msgBytes)
	}
}
func HandleMessage(hub *hubs.MainHub, grpcClient *clients.MessageServiceClient, clientModel models.Client, msg []byte) {
	log.Println("Handle Message called.")
	var input models.IncomingMessage
	if err := json.Unmarshal(msg, &input); err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}
	// Switch type
	switch input.Type {
	case models.SendMessage:
		log.Println("SendMessage type")
		// Send message through WebSocket
		chatId := clientModel.Chat.ID
		userId := input.UserId
		message := models.NewMessage(input.Text, userId, chatId)
		err := hub.SendMessageToClientWithModel(userId, *message)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		// Send message to gRPC service
		err = grpcClient.SendMessage(context.Background(), chatId.String(), userId.String(), input.Text)
		if err != nil {
			log.Println("Error sending message to gRPC service:", err)
			return
		}
		log.Printf("Message sent to chat %s by user %s: %s", chatId, userId, input.Text)

	case models.ConnectUserToChat:
		log.Println("ConnectUserToChat type")
		chatId, err := uuid.Parse(input.ChatId.String())
		if err != nil {
			log.Println("Error parsing chatId:", err)
			return
		}
		hub.ConnectUserToChat(clientModel.UserId, chatId)
		log.Printf("User %s connected to chat %s", clientModel.UserId, chatId)
	case models.DisconnectUserFromChat:
		log.Println("DisconnectUserFromChat type")
		chatId, err := uuid.Parse(input.ChatId.String())
		if err != nil {
			log.Println("Error parsing chatId:", err)
			return
		}
		hub.DisconnectUserFromChat(clientModel.UserId)
		log.Printf("User %s disconnected from chat %s", clientModel.UserId, chatId)
	}
}
