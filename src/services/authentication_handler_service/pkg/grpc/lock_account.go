package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

// LockAccount locks an account as long as it exists from the context of the authentication service
func (*Server) LockAccount(ctx context.Context, req *proto.LockAccountRequest) (*proto.LockAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - lock api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method LockAccount not implemented")
}
