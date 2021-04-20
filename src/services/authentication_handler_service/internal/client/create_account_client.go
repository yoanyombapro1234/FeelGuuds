package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// CreateAccount creates an account
func (c *Client) CreateAccount(ctx context.Context, email, password string) (uint32, error) {
	request := &proto.CreateAccountRequest{
		Email: email,
		Password: password,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	response, err := c.conn.CreateAccount(ctx, request)
	if err != nil {
		if er, ok := status.FromError(err); ok {
			return 0, fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return 0, fmt.Errorf("server: %s", err.Error())
	}

	return response.Id, nil
}
