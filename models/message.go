package models

import (
	"github.com/google/uuid"
)

type Message struct {
	Text   string    `json:"text"`
	UserId uuid.UUID `json:"userId"`
	ChatId uuid.UUID `json:"chatId"`
}

func NewMessage(text string, userId, chatId uuid.UUID) *Message {
	return &Message{
		Text:   text,
		UserId: userId,
		ChatId: chatId,
	}
}
