package main

import (
	"context"
	"errors"
	wallet_app "github.com/rima971/wallet-app/authenticator"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":8090"

type AuthenticatorServerImpl struct {
	wallet_app.UnimplementedAuthenticatorServer
}

func (s AuthenticatorServerImpl) Register(ctx context.Context, req *wallet_app.User) (*wallet_app.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("bad request: username or password is empty")
	}
	res := &wallet_app.RegisterResponse{
		Message: "user registered successfully",
		User:    req,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &AuthenticatorServerImpl{}

	wallet_app.RegisterAuthenticatorServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %v", err)
	}
}
