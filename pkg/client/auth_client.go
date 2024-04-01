package client

import (
	"context"
	"flotta-home/mindbond/websocket-server/pkg/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
}

func InitAuthServiceClient(url string) AuthServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	return AuthServiceClient{
		Client: pb.NewAuthServiceClient(cc),
	}
}

func (a *AuthServiceClient) Validate(token string) (*pb.ValidateResponse, error) {
	request := &pb.ValidateRequest{
		Token: token,
	}
	return a.Client.Validate(context.Background(), request)
}
