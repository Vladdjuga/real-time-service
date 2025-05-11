package clients

import (
	"github.com/gin-gonic/gin"
	"realTimeService/gen/chatpb"
)

type ChatServiceClient struct {
	client chatpb.ChatServiceClient
}

func NewChatServiceClient(client chatpb.ChatServiceClient) *ChatServiceClient {
	return &ChatServiceClient{
		client: client,
	}
}
func (c *ChatServiceClient) UserChatExists(ctx *gin.Context, userId, chatId string) (bool, error) {
	resp, err := c.client.UserChatExists(ctx, &chatpb.UserChatExistsRequest{
		UserId: userId,
		ChatId: chatId,
	})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}
