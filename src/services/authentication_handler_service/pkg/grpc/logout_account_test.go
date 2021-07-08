package grpc

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

func Test_logout_account(t *testing.T) {
	id := "1"

	tests := []struct {
		scenario          string
		id                string
		res               *proto.LogoutAccountResponse
		errCode           codes.Code
		errMsg            string
		LogoutAccountFunc func() error
	}{
		// scenario: unlock account that exists
		{
			"log out valid existing account",
			id,
			&proto.LogoutAccountResponse{},
			codes.Unknown,
			"",
			func() error {
				return nil
			},
		},
		// scenario: failed to logout account
		{
			"account already locked",
			id,
			&proto.LogoutAccountResponse{},
			codes.Unknown,
			service_errors.ErrCannotLogoutAccount.Error(),
			func() error {
				return service_errors.ErrCannotLogoutAccount
			},
		},
		// scenario: invalid request
		{
			"invalid request",
			"0",
			nil,
			codes.InvalidArgument,
			"invalid input argument",
			func() error {
				return service_errors.ErrInvalidInputArguments
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.LogOutAccountFunc = tt.LogoutAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			accountId, _ := strconv.Atoi(tt.id)
			request := &proto.LogoutAccountRequest{
				Id: uint32(accountId),
			}

			response, err := client.LogoutAccount(ctx, request)

			if response != nil {
				if response.GetError() != tt.res.GetError() {
					t.Error("response: expected", tt.res.GetError(), "received", response.GetError())
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
					if !strings.Contains(er.Message(), tt.errMsg) {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
