package clients

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"realTimeService/gen/messengerpb"
	"time"
)

type MessageServiceClient struct {
	client messengerpb.MessengerClient
}

func NewMessageServiceClient(conn *grpc.ClientConn) *MessageServiceClient {
	return &MessageServiceClient{
		client: messengerpb.NewMessengerClient(conn),
	}
}

func (m *MessageServiceClient) SendMessage(ctx context.Context, chatId, userId string, text, authToken string) error {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + authToken,
	})
	ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)
	log.Println("Sending message: ", chatId, userId, text)
	_, err := m.client.SendMessage(ctxWithMetadata,
		&messengerpb.SendMessageRequest{
			ChatId: chatId,
			UserId: userId,
			Text:   text,
			SentAt: timestamppb.New(time.Now().UTC()),
		},
	)
	if err != nil {
		log.Fatalln("Error sending message: ", err)
		return err
	}
	log.Println("Message sent successfully")
	return nil
}
func (m *MessageServiceClient) GetMessage(ctx context.Context, chatId, userId, authToken string) (*messengerpb.GetMessageResponse, error) {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + authToken,
	})
	ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)
	log.Println("Getting message: ", chatId, userId)
	resp, err := m.client.GetMessage(ctxWithMetadata,
		&messengerpb.GetMessageRequest{
			ChatId: chatId,
			UserId: userId,
		},
	)
	if err != nil {
		log.Fatalln("Error getting message: ", err)
		return nil, err
	}
	log.Println("Message received successfully: ", resp)
	return resp, nil
}
