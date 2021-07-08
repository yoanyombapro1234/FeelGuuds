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

func Test_unlock_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	id := "1"

	tests := []struct {
		scenario          string
		id                string
		res               *proto.UnLockAccountResponse
		errCode           codes.Code
		errMsg            string
		UnLockAccountFunc func(id string) error
	}{
		// scenario: unlock account that exists
		{
			"lock account that exists",
			id,
			&proto.UnLockAccountResponse{},
			codes.Unknown,
			"",
			func(id string) error {
				return nil
			},
		},
		// scenario: failed to unlock account - account already unlocked
		{
			"account already locked",
			id,
			&proto.UnLockAccountResponse{},
			codes.Unknown,
			service_errors.ErrAccountAlreadyUnlocked.Error(),
			func(id string) error {
				return service_errors.ErrAccountAlreadyUnlocked
			},
		},
		// scenario: failed to unlock account - account does not exist
		{
			"account already exists",
			id,
			&proto.UnLockAccountResponse{},
			codes.Unknown,
			service_errors.ErrAccountDoesNotExist.Error(),
			func(id string) error {
				return service_errors.ErrAccountDoesNotExist
			},
		},
		// scenario: invalid request
		{
			"invalid request",
			"0",
			nil,
			codes.InvalidArgument,
			"invalid input arguments",
			func(id string) error {
				return service_errors.ErrInvalidInputArguments
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.UnlockAccountFunc = tt.UnLockAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			accountId, _ := strconv.Atoi(tt.id)
			request := &proto.UnLockAccountRequest{
				Id: uint32(accountId),
			}

			response, err := client.UnLockAccount(ctx, request)

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
