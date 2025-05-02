package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Content   string    `json:"content"`
	UserId    uuid.UUID `json:"userId"`
	ChatId    uuid.UUID `json:"chatId"`
	Timestamp time.Time `json:"timestamp"`
}
