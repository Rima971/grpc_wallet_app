package main

import (
	"context"
	"errors"
	wallet_app "github.com/rima971/wallet-app/authenticator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func server(ctx context.Context) (wallet_app.AuthenticatorClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	wallet_app.RegisterAuthenticatorServer(baseServer, AuthenticatorServerImpl{})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %s", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %s", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := wallet_app.NewAuthenticatorClient(conn)

	return client, closer
}

func TestAuthenticatorServer_Register(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *wallet_app.RegisterResponse
		err error
	}
	tests := map[string]struct {
		in       *wallet_app.User
		expected expectation
	}{
		"success": {
			in: &wallet_app.User{
				Username: "user",
				Password: "password",
			},
			expected: expectation{
				out: &wallet_app.RegisterResponse{
					Message: "user registered successfully",
					User: &wallet_app.User{
						Username: "user",
						Password: "password",
					},
				},
				err: nil,
			},
		},
		"bad request: empty username": {
			in: &wallet_app.User{
				Username: "",
				Password: "password",
			},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = bad request: username or password is empty"),
			},
		},
		"bad request: empty password": {
			in: &wallet_app.User{
				Username: "user",
				Password: "",
			},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = bad request: username or password is empty"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Register(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Message != out.Message ||
					tt.expected.out.User.Username != out.User.Username ||
					tt.expected.out.User.Password != out.User.Password {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}

		})
	}
}
