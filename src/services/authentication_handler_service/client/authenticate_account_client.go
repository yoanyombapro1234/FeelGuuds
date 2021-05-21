package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// CreateAccount creates an account
func (c *Client) AuthenticateAccount(ctx context.Context, email, password string) (string, error) {
	request := &proto.AuthenticateAccountRequest{
		Email:    email,
		Password: password,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	response, err := c.conn.AuthenticateAccount(ctx, request)
	if err != nil {
		if er, ok := status.FromError(err); ok {
			return "", fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return "", fmt.Errorf("server: %s", err.Error())
	}

	return response.Token, nil
}
