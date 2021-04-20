package client

import (
	"time"

	"google.golang.org/grpc"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

type Client struct {
	conn proto.AuthenticationHandlerServiceApiClient
	timeout time.Duration
}

// NewClient instantiates a new instance of the service's client
func NewClient(conn *grpc.ClientConn, timeout time.Duration) *Client {
	return &Client{
		conn:    proto.NewAuthenticationHandlerServiceApiClient(conn),
		timeout: timeout,
	}
}



