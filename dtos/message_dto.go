package dtos

import (
	"github.com/google/uuid"
	"time"
)

type MessageStatus int

const (
	Pending MessageStatus = 1 << iota
	Sent
	Delivered
	Read
)

type MessageDto struct {
	ID         uuid.UUID     `json:"id"`
	SentAt     time.Time     `json:"sentAt"`
	ReceivedAt time.Time     `json:"receivedAt"`
	Text       string        `json:"text"`
	UserID     uuid.UUID     `json:"userId"`
	ChatID     uuid.UUID     `json:"chatId"`
	Status     MessageStatus `json:"status"`
}

func NewMessageDto(
	id, userId, chatId uuid.UUID,
	text string,
	sentAt, receivedAt time.Time,
	status MessageStatus) *MessageDto {
	return &MessageDto{
		ID:         id,
		SentAt:     sentAt,
		ReceivedAt: receivedAt,
		Text:       text,
		UserID:     userId,
		ChatID:     chatId,
		Status:     status,
	}
}
