package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// LogoutAccount revokes the user account session from the context of the authentication service
func (*Server) LogoutAccount(ctx context.Context, req *proto.LogoutAccountRequest) (*proto.LogoutAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		AUTHENTICATED
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - update api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method LogoutAccount not implemented")
}
