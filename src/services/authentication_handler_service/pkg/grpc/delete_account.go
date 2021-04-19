package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// DeleteAccount deletes a user account via the authentication service
func (*Server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		AUTHENTICATED
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - delete api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccount not implemented")
}
