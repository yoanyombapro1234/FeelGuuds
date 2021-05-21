package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// GetAccount obtains an account from the authentication service by account Id
func (c *Client) GetAccount(ctx context.Context, accountId uint32) (*proto.Account, error) {
	request := &proto.GetAccountRequest{
		Id: accountId,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	response, err := c.conn.GetAccount(ctx, request)
	if err != nil {
		if er, ok := status.FromError(err); ok {
			return nil, fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return nil, fmt.Errorf("server: %s", err.Error())
	}

	return response.Account, nil
}
