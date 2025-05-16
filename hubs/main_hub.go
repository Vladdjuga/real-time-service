package hubs

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"realTimeService/models"
	"sync"
)

type MainHub struct {
	Clients map[uuid.UUID]*models.Client
	Chats   map[uuid.UUID]*models.Chat
	mut     sync.RWMutex
}

func NewMainHub() *MainHub {
	return &MainHub{
		Clients: make(map[uuid.UUID]*models.Client),
		Chats:   make(map[uuid.UUID]*models.Chat),
		mut:     sync.RWMutex{},
	}
}
func (h *MainHub) AddClient(client *models.Client) {
	h.mut.Lock()
	defer h.mut.Unlock()
	h.Clients[client.UserId] = client
}

// ConnectUserToChat This function is used to add a chatId to a client.
// It is called when a user joins a chat.
// This function will later be removed.
func (h *MainHub) ConnectUserToChat(userId, chatId uuid.UUID) {
	h.mut.Lock()
	defer h.mut.Unlock()
	client, ok := h.Clients[userId]
	if !ok {
		return
	}
	chat, ok := h.Chats[chatId]
	if !ok {
		chat = models.NewChat(chatId)
		h.Chats[chatId] = chat
	}
	client.Chat = chat
	chat.AddClient(client)
}

// DisconnectUserFromChat This function is used to remove a chatId from a client.
// It is called when a user leaves a chat.
// This function will later be removed.
func (h *MainHub) DisconnectUserFromChat(userId uuid.UUID) {
	h.mut.Lock()
	defer h.mut.Unlock()
	client, ok := h.Clients[userId]
	if !ok {
		return
	}
	chat, ok := h.Chats[client.Chat.ID]
	if !ok {
		return
	}
	client.Chat = nil
	chat.RemoveClient(client)
}

func (h *MainHub) RemoveClient(userId uuid.UUID) {
	h.mut.Lock()
	defer h.mut.Unlock()
	delete(h.Clients, userId)
}
func (h *MainHub) GetClient(userId uuid.UUID) *models.Client {
	h.mut.RLock()
	defer h.mut.RUnlock()
	client, ok := h.Clients[userId]
	if !ok {
		return nil
	}
	return client
}
func (h *MainHub) SendMessageToChat(chatId uuid.UUID, message models.Message) error {
	h.mut.RLock()
	defer h.mut.RUnlock()
	chat, ok := h.Chats[chatId]
	if !ok {
		return fmt.Errorf("chat not found")
	}
	for _, client := range chat.Clients {
		err := h.sendMessageToClientWithModel(client, message)
		if err != nil {
			err = fmt.Errorf("error sending message to client: %w", err)
			log.Println(err)
			continue
		}
	}
	return nil
}
func (h *MainHub) sendMessageToClient(client *models.Client, message []byte) error {
	h.mut.RLock()
	defer h.mut.RUnlock()
	log.Println("Sending message to client : ", client.UserId)
	return client.Conn.WriteMessage(websocket.TextMessage, message)
}
func (h *MainHub) sendMessageToClientWithModel(client *models.Client, message models.Message) error {
	h.mut.RLock()
	defer h.mut.RUnlock()
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	log.Println("Sending message to client : ", client.UserId)
	return client.Conn.WriteMessage(websocket.TextMessage, messageBytes)
}
