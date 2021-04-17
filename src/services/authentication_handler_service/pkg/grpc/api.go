package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

func (*Server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (*Server) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
func (*Server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccount not implemented")
}
func (*Server) LockAccount(ctx context.Context, req *proto.LockAccountRequest) (*proto.LockAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LockAccount not implemented")
}
func (*Server) UnLockAccount(ctx context.Context, req *proto.UnLockAccountRequest) (*proto.UnLockAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnLockAccount not implemented")
}
func (*Server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (*Server) AuthenticateAccount(ctx context.Context, req *proto.AuthenticateAccountRequest) (*proto.AuthenticateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateAccount not implemented")
}
func (*Server) LogoutAccount(ctx context.Context, req *proto.LogoutAccountRequest) (*proto.LogoutAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogoutAccount not implemented")
}
