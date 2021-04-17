package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

func (*Server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		1. Decode the request body, validate request params, and emit any necessary metrics
		2. Through RPC invoke the authentication service (again using timeout + retry logic, proper instrumentation, and traces)
			- On request failure, we check the error status return the error
	*/

	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (*Server) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		AUTHENTICATED
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - update api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
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
func (*Server) UnLockAccount(ctx context.Context, req *proto.UnLockAccountRequest) (*proto.UnLockAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - unlock api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method UnLockAccount not implemented")
}
func (*Server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		AUTHENTICATED
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - update api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (*Server) AuthenticateAccount(ctx context.Context, req *proto.AuthenticateAccountRequest) (*proto.AuthenticateAccountResponse, error) {
	// TODO: start span + emit any preliminary metrics
	// TODO: start operation as timeout op
	/*
		1. Decode the request body, validate request parameters,
		2. Through RPC invoke the authentication service - login api (timeouts, retry logic, instrumentation, traces, ...etc)
			- return response
	*/
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateAccount not implemented")
}
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
