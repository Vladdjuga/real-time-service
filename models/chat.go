package models

import "github.com/google/uuid"

type Chat struct {
	ID       uuid.UUID `json:"id"`
	Messages []Message `json:"messages"`
}

func NewChat(id uuid.UUID) *Chat {
	return &Chat{
		ID:       id,
		Messages: []Message{},
	}
}

func (chat *Chat) AddMessage(message Message) {
	chat.Messages = append(chat.Messages, message)
}
