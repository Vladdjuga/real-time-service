package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"realTimeService/interfaces"
	"realTimeService/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WsHandler struct {
	container interfaces.Container
}

func NewWsHandler(c interfaces.Container) *WsHandler {
	return &WsHandler{
		container: c,
	}
}

// Handler interface implementation
func (h *WsHandler) Handle(ctx *gin.Context) {
	log.Println("WsHandler called.")
	// This userId is passed from the AuthMiddleware
	// and is used to identify the user in the WebSocket connection.
	userId := ctx.GetString("user_sub")
	token := ctx.GetString("auth_token")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	var client *models.Client
	if err != nil {
		log.Println("Write error:", err)
		return
	} else {
		log.Println("Connected to websocket")
		client = models.NewClient(uuid.MustParse(userId), nil, conn)
		h.container.GetHub().AddClient(client)
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
		err = h.HandleMessage(ctx, client, token, msgBytes)
		if err != nil {
			log.Println("Error handling message:", err)
			break
		}
	}
}
func (h *WsHandler) HandleMessage(ctx *gin.Context, clientModel *models.Client, token string, msg []byte) error {
	log.Println("Handle Message called.")
	var input models.IncomingMessage
	if err := json.Unmarshal(msg, &input); err != nil {
		log.Println("Error unmarshalling message:", err)
		return err
	}
	// Switch type
	switch input.Type {
	case models.SendMessage:
		if err := h.sendMessage(ctx, clientModel, input, token); err != nil {
			log.Println("Error sending message:", err)
			return err
		}
	case models.ConnectUserToChat:
		chatId, err := h.connectUserToChat(ctx, input, clientModel, token)
		if err != nil {
			log.Println("Error connecting user to chat:", err)
			return err
		}
		log.Printf("User %s connected to chat %s", clientModel.UserId, chatId)
	case models.DisconnectUserFromChat:
		if err := h.disconnectUserFromChat(input, clientModel); err != nil {
			log.Println("Error disconnecting user from chat:", err)
			return err
		}
	}
	log.Println("Message handled successfully")
	return nil
}
func (h *WsHandler) sendMessage(ctx *gin.Context, clientModel *models.Client, input models.IncomingMessage, token string) error {
	log.Println("SendMessage type")
	// Check if the user is connected to a chat
	// If not, connect the user to the chat
	if clientModel.Chat == nil {
		_, err := h.connectUserToChat(ctx, input, clientModel, token)
		if err != nil {
			log.Println("Error connecting user to chat:", err)
			return err
		}
	}
	// Send message through WebSocket
	chatId := clientModel.Chat.ID
	message := models.NewMessage(input.Text, clientModel.UserId, chatId)
	err := h.container.GetHub().SendMessageToChat(chatId, *message)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	// Send message to gRPC service
	err = h.container.GetMessageClient().SendMessage(context.Background(),
		chatId.String(), clientModel.UserId.String(), input.Text, token)
	if err != nil {
		log.Println("Error sending message to gRPC service:", err)
		return err
	}
	log.Printf("Message sent to chat %s by user %s: %s", chatId, clientModel.UserId, input.Text)
	return nil
}

func (h *WsHandler) disconnectUserFromChat(input models.IncomingMessage, clientModel *models.Client) error {
	log.Println("DisconnectUserFromChat type")
	chatId, err := uuid.Parse(input.ChatId.String())
	if err != nil {
		log.Println("Error parsing chatId:", err)
		return err
	}
	h.container.GetHub().DisconnectUserFromChat(clientModel.UserId)
	log.Printf("User %s disconnected from chat %s", clientModel.UserId, chatId)
	return nil
}

func (h *WsHandler) connectUserToChat(ctx *gin.Context, input models.IncomingMessage, clientModel *models.Client, token string) (uuid.UUID, error) {
	log.Println("ConnectUserToChat type")
	res, err := h.container.GetChatClient().UserChatExists(ctx, clientModel.UserId.String(), input.ChatId.String(), token)
	if err != nil {
		log.Println("Error checking if user is in chat:", err)
		return uuid.UUID{}, err
	}
	if !res {
		log.Println("User is not in chat")
		err := fmt.Errorf("user is not in chat")
		return uuid.UUID{}, err
	}
	chatId, err := uuid.Parse(input.ChatId.String())
	if err != nil {
		log.Println("Error parsing chatId:", err)
		return uuid.UUID{}, err
	}
	h.container.GetHub().ConnectUserToChat(clientModel.UserId, chatId)
	return chatId, nil
}
