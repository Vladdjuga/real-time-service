package clients

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"realTimeService/gen/chatpb"
)

type ChatServiceClient struct {
	client chatpb.ChatServiceClient
}

func NewChatServiceClient(conn *grpc.ClientConn) *ChatServiceClient {
	return &ChatServiceClient{
		client: chatpb.NewChatServiceClient(conn),
	}
}
func (c *ChatServiceClient) UserChatExists(ctx *gin.Context, userId, chatId, authToken string) (bool, error) {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + authToken,
	})
	ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)
	log.Println("Checking if user chat exists")
	resp, err := c.client.UserChatExists(ctxWithMetadata, &chatpb.UserChatExistsRequest{
		UserId: userId,
		ChatId: chatId,
	})
	if err != nil {

		log.Fatalln("Error calling UserChatExists: ", err)
		return false, err
	}
	log.Println("User chat exists: ", resp)
	return resp.Exists, nil
}
