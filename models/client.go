package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserId uuid.UUID
	Chat   *Chat
	Conn   *websocket.Conn
}

func NewClient(userId uuid.UUID, chat *Chat, conn *websocket.Conn) *Client {
	return &Client{
		UserId: userId,
		Chat:   chat,
		Conn:   conn,
	}
}
