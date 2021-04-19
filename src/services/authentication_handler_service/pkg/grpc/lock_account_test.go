package grpc

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

func Test_lock_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	expectedErrMsg := "retry limit reached"

	id := "1"

	tests := []struct {
		scenario          string
		id          string
		res               *proto.LockAccountResponse
		errCode           codes.Code
		errMsg            string
		LockAccountFunc func(id string) (error)
	}{
		// scenario: lock account that exists
		{
			"lock account that exists",
			id,
			&proto.LockAccountResponse{},
			codes.Unknown,
			"",
			func(id string) error {
				return nil
			},
		},
		// scenario: failed to lock account - account already locked
		{
			"account already locked",
			id,
			&proto.LockAccountResponse{},
			codes.Unknown,
			expectedErrMsg,
			func(id string) error {
				return errors.New("account already locked")
			},
		},
		// scenario: failed to lock account - account does not exist
		{
			"account already exists",
			id,
			&proto.LockAccountResponse{},
			codes.Unknown,
			expectedErrMsg,
			func(id string) error {
				return errors.New("account does not exist")
			},
		},
		// scenario: invalid request
		{
			"invalid request",
			"0",
			nil,
			codes.InvalidArgument,
			"invalid input argument",
			func(id string) error {
				return service_errors.ErrInvalidInputArguments
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.LockAccountFunc = tt.LockAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			accountId , _ := strconv.Atoi(tt.id)
			request := &proto.LockAccountRequest{
				Id:    uint32(accountId),
			}

			response, err := client.LockAccount(ctx, request)

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
