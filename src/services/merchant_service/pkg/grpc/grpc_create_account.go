package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
)

func (s *Server) CreateAccount(ctx context.Context, request *merchant_service_proto_v1.CreateAccountRequest) (*merchant_service_proto_v1.CreateAccountResponse, error) {
	// 1. perform request validations
	// 2. create account record from context of auth service
	// 3. invoke and create an account from the context of stripe
	// 4. store record in database after applying mutations
	// 5. log in the account by invoking the auth service (authenticate api)

	panic("implement me")
}
