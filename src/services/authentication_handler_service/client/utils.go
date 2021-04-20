package client

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	s_grpc "github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/grpc"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	s := s_grpc.NewMockServer(nil)
	proto.RegisterAuthenticationHandlerServiceApiServer(server, s)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

// createAccountTestHelper creates a user account based on the passed in credentials. This should only be used in unit tests
func createAccountTestHelper(c *Client, ctx context.Context, email, password string) (
	uint32, error) {
	id, err := c.CreateAccount(ctx, email, password)
	return id, err
}

