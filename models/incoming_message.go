package models

import "github.com/google/uuid"

type Type string

const (
	SendMessage            Type = "sendMessage"
	ConnectUserToChat           = Type("connectUserToChat")
	DisconnectUserFromChat      = Type("disconnectUserFromChat")
)

type IncomingMessage struct {
	Type   Type      `json:"type"`
	UserId uuid.UUID `json:"user_id"`
	ChatId uuid.UUID `json:"chat_id"`
	Text   string    `json:"text"`
}
