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
	// This userId is passed from the AuthMiddleware
	// and is used to identify the user in the WebSocket connection.
	userId := c.GetString("user_sub")
	token := c.GetString("auth_token")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	var client *models.Client
	if err != nil {
		log.Println("Write error:", err)
		return
	} else {
		log.Println("Connected to websocket")
		client = models.NewClient(uuid.MustParse(userId), nil, conn)
		hub.AddClient(client)
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
		HandleMessage(hub, grpcClient, client, token, msgBytes)
	}
}
func HandleMessage(hub *hubs.MainHub, grpcClient *clients.MessageServiceClient, clientModel *models.Client, token string, msg []byte) {
	log.Println("Handle Message called.")
	var input models.IncomingMessage
	if err := json.Unmarshal(msg, &input); err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}
	// Switch type
	switch input.Type {
	case models.SendMessage:
		if err := sendMessage(hub, grpcClient, clientModel, input, token); err != nil {
			log.Println("Error sending message:", err)
			return
		}
	case models.ConnectUserToChat:
		chatId, err := connectUserToChat(hub, input, clientModel)
		if err != nil {
			log.Println("Error connecting user to chat:", err)
			return
		}
		log.Printf("User %s connected to chat %s", clientModel.UserId, chatId)
	case models.DisconnectUserFromChat:
		if err := disconnectUserFromChat(hub, input, clientModel); err != nil {
			log.Println("Error disconnecting user from chat:", err)
			return
		}
	}
}
func sendMessage(hub *hubs.MainHub, grpcClient *clients.MessageServiceClient,
	clientModel *models.Client, input models.IncomingMessage, token string) error {
	log.Println("SendMessage type")
	// Check if the user is connected to a chat
	// If not, connect the user to the chat
	if clientModel.Chat == nil {
		_, err := connectUserToChat(hub, input, clientModel)
		if err != nil {
			log.Println("Error connecting user to chat:", err)
			return err
		}
	}
	// Send message through WebSocket
	chatId := clientModel.Chat.ID
	message := models.NewMessage(input.Text, clientModel.UserId, chatId)
	err := hub.SendMessageToChat(chatId, *message)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	// Send message to gRPC service
	err = grpcClient.SendMessage(context.Background(), chatId.String(), clientModel.UserId.String(), input.Text, token)
	if err != nil {
		log.Println("Error sending message to gRPC service:", err)
		return err
	}
	log.Printf("Message sent to chat %s by user %s: %s", chatId, clientModel.UserId, input.Text)
	return nil
}

func disconnectUserFromChat(hub *hubs.MainHub, input models.IncomingMessage, clientModel *models.Client) error {
	log.Println("DisconnectUserFromChat type")
	chatId, err := uuid.Parse(input.ChatId.String())
	if err != nil {
		log.Println("Error parsing chatId:", err)
		return err
	}
	hub.DisconnectUserFromChat(clientModel.UserId)
	log.Printf("User %s disconnected from chat %s", clientModel.UserId, chatId)
	return nil
}

func connectUserToChat(hub *hubs.MainHub, input models.IncomingMessage, clientModel *models.Client) (uuid.UUID, error) {
	log.Println("ConnectUserToChat type")
	chatId, err := uuid.Parse(input.ChatId.String())
	if err != nil {
		log.Println("Error parsing chatId:", err)
		return uuid.UUID{}, err
	}
	hub.ConnectUserToChat(clientModel.UserId, chatId)
	return chatId, nil
}
