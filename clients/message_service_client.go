package clients

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	_, err := m.client.SendMessage(ctxWithMetadata,
		&messengerpb.SendMessageRequest{
			ChatId: chatId,
			UserId: userId,
			Text:   text,
			SentAt: timestamppb.New(time.Now().UTC()),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (m *MessageServiceClient) GetMessage(ctx context.Context, chatId, userId, authToken string) (*messengerpb.GetMessageResponse, error) {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + authToken,
	})
	ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)
	resp, err := m.client.GetMessage(ctxWithMetadata,
		&messengerpb.GetMessageRequest{
			ChatId: chatId,
			UserId: userId,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
