package models

import "github.com/google/uuid"

type Chat struct {
	ID      uuid.UUID `json:"id"`
	Clients []*Client `json:"Clients"`
}

func NewChat(id uuid.UUID) *Chat {
	return &Chat{
		ID:      id,
		Clients: []*Client{},
	}
}

func (chat *Chat) AddClient(client *Client) {
	chat.Clients = append(chat.Clients, client)
}
func (chat *Chat) filterClients(shouldKeep func(*Client) bool) {
	var result []*Client
	for _, client := range chat.Clients {
		if shouldKeep(client) {
			result = append(result, client)
		}
	}
	chat.Clients = result
}
func (chat *Chat) RemoveClient(client *Client) {
	chat.filterClients(func(c *Client) bool {
		return c.UserId != client.UserId
	})
}
