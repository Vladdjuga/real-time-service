package interfaces

import (
	"realTimeService/clients"
	"realTimeService/configuration"
	"realTimeService/hubs"
)

type Container interface {
	GetHub() *hubs.MainHub
	GetMessageClient() *clients.MessageServiceClient
	GetChatClient() *clients.ChatServiceClient
	InitializeProviders(cfg *configuration.Config)
	Close() error
}
