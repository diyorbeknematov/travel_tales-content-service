package server

import (
	"content-service/config"
	"content-service/generated/user"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(cfg config.Config) (user.AuthServiceClient, error) {
	connUser, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", cfg.USER_CLIENT_PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	userClient := user.NewAuthServiceClient(connUser)
	return userClient, nil
}