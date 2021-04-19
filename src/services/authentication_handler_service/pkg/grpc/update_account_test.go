package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
)

func Test_update_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	expectedErrMsg := "retry limit reached"

	email := fmt.Sprintf("test_%s@gmail.com", GenerateRandomString(17))

	tests := []struct {
		scenario          string
		email             string
		id          int
		res               *proto.UpdateAccountResponse
		errCode           codes.Code
		errMsg            string
		UpdateAccoutnFunc func(id, email string) error
	}{
		// scenario: successful update of account
		{
			"account exists based on credentials",
			email,
			1,
			nil,
			codes.Unknown,
			"",
			func(id, username string) error {
				return nil
			},
		},
		// scenario: invalid email params
		{
			"invalid email params",
			"",
			0,
			nil,
			codes.InvalidArgument,
			"invalid input arguments",
			func(id, username string) error {
				return nil
			},
		},
		// scenario: account does not exist
		{
			"account does not exist",
			email,
			1,
			nil,
			codes.Unknown,
			expectedErrMsg,
			func(id, username string) error {
				return errors.New("account does not exit")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.UpdateAccountFunc = tt.UpdateAccoutnFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			request := &proto.UpdateAccountRequest{
				Email:    tt.email,
				Id: uint32(tt.id),
			}

			response, err := client.UpdateAccount(ctx, request)

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
