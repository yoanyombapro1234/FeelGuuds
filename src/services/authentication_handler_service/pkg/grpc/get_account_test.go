package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/util"
)

func Test_get_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_password_%s", util.GenerateRandomString(17))
	testAccount := &proto.Account{
		Id:       1,
		Username: email,
		Locked:   true,
		Deleted:  true,
	}

	tests := []struct {
		scenario       string
		id             uint32
		email          string
		password       string
		res            *proto.GetAccountResponse
		errCode        codes.Code
		errMsg         string
		GetAccountFunc func(id string) (*core_auth_sdk.Account, error)
	}{
		// scenario: valid request
		{
			"account successfully retrieved",
			1,
			email,
			password,
			&proto.GetAccountResponse{
				Account: testAccount,
				Error:   "",
			},
			codes.Unknown,
			"",
			func(id string) (*core_auth_sdk.Account, error) {
				return &core_auth_sdk.Account{
					ID:       int(testAccount.Id),
					Username: testAccount.Username,
				}, nil
			},
		},
		// scenario: invalid request
		{
			"invalid request - invalid input arguments",
			0,
			email,
			"",
			nil,
			codes.InvalidArgument,
			"invalid input arguments",
			func(id string) (*core_auth_sdk.Account, error) {
				return nil, service_errors.ErrInvalidInputArguments
			},
		},
		// scenario: invalid request
		{
			"invalid request - account does not exist",
			1,
			email,
			"",
			nil,
			codes.Unknown,
			"account does not exist",
			func(id string) (*core_auth_sdk.Account, error) {
				return nil, errors.New("account does not exist")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.GetAccountFunc = tt.GetAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			request := &proto.GetAccountRequest{
				Id: tt.id,
			}

			response, err := client.GetAccount(ctx, request)

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
