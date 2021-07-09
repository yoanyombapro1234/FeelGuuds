package grpc_test_utils

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
)

type MockedService struct{}

func (s *MockedService) UpdateAccount(ctx context.Context, request *merchant_service_proto_v1.UpdateAccountRequest) (*merchant_service_proto_v1.UpdateAccountResponse, error) {
	panic("implement me")
}

func (s *MockedService) DeleteAccount(ctx context.Context, request *merchant_service_proto_v1.DeleteAccountRequest) (*merchant_service_proto_v1.DeleteAccountResponse, error) {
	panic("implement me")
}

func (s *MockedService) GetAccount(ctx context.Context, request *merchant_service_proto_v1.GetAccountRequest) (*merchant_service_proto_v1.GetAccountResponse, error) {
	panic("implement me")
}

func (s *MockedService) GetAccounts(ctx context.Context, request *merchant_service_proto_v1.GetAccountsRequest) (*merchant_service_proto_v1.GetAccountsResponse, error) {
	panic("implement me")
}

func (s *MockedService) SetAccountStatus(ctx context.Context, request *merchant_service_proto_v1.SetAccountStatusRequest) (*merchant_service_proto_v1.SetAccountStatusResponse, error) {
	panic("implement me")
}

func (s *MockedService) StartAccountOnboarding(ctx context.Context, request *merchant_service_proto_v1.StartAccountOnboardingRequest) (*merchant_service_proto_v1.StartAccountOnboardingRespone, error) {
	panic("implement me")
}

var _ merchant_service_proto_v1.MerchantServiceServer = (*MockedService)(nil)

func (s *MockedService) CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*merchant_service_proto_v1.CreateAccountResponse, error) {
	return &merchant_service_proto_v1.CreateAccountResponse{AccountId: uint64(1000)}, nil
}
