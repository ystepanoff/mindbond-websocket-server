package client

import (
	"context"
	"flotta-home/mindbond/websocket-server/pkg/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatServiceClient struct {
	Client pb.ChatServiceClient
}

func InitChatServiceClient(url string) ChatServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	return ChatServiceClient{
		Client: pb.NewChatServiceClient(cc),
	}
}

func (a *ChatServiceClient) AddMessage(userFromId, userToId int64, message, token string) (*pb.AddMessageResponse, error) {
	request := &pb.AddMessageRequest{
		UserFromId: userFromId,
		UserToId:   userToId,
		Message:    message,
		Token:      token,
	}
	return a.Client.AddMessage(context.Background(), request)
}
