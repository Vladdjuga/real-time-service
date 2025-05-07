package clients

import (
	"context"
	"google.golang.org/grpc"
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

func (m *MessageServiceClient) SendMessage(ctx context.Context, chatId, userId string, text string) error {
	_, err := m.client.SendMessage(ctx,
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
func (m *MessageServiceClient) GetMessage(ctx context.Context, chatId, userId string) (*messengerpb.GetMessageResponse, error) {
	resp, err := m.client.GetMessage(ctx,
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
