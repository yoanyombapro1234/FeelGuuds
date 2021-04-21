package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// LockAccount locks an account
func (c *Client) LockAccount(ctx context.Context, accountId uint32) error {
	request := &proto.LockAccountRequest{
		Id: accountId,
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	_, err := c.conn.LockAccount(ctx, request)
	if err != nil {
		if er, ok := status.FromError(err); ok {
			return fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return fmt.Errorf("server: %s", err.Error())
	}

	return nil
}