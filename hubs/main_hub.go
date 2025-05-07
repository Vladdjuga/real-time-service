package hubs

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"realTimeService/models"
	"sync"
)

type MainHub struct {
	Clients map[uuid.UUID]*models.Client
	mut     sync.RWMutex
}

func NewMainHub() *MainHub {
	return &MainHub{
		Clients: make(map[uuid.UUID]*models.Client),
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
	client.Chat = models.NewChat(chatId)
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
	client.Chat = nil
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
func (h *MainHub) SendMessageToClient(userId uuid.UUID, message []byte) error {
	h.mut.RLock()
	defer h.mut.RUnlock()
	client, ok := h.Clients[userId]
	if !ok {
		return fmt.Errorf("user not connected")
	}
	return client.Conn.WriteMessage(websocket.TextMessage, message)
}
func (h *MainHub) SendMessageToClientWithModel(userId uuid.UUID, message models.Message) error {
	h.mut.RLock()
	defer h.mut.RUnlock()
	client, ok := h.Clients[userId]
	if !ok {
		return fmt.Errorf("user not connected")
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	return client.Conn.WriteMessage(websocket.TextMessage, messageBytes)
}
