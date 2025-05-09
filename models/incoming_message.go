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
	ChatId uuid.UUID `json:"chatId"`
	Text   string    `json:"text"`
}
