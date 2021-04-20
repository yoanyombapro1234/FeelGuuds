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

func Test_create_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	expectedErrMsg := "retry limit reached"
	ThirdPartyMockService.ImportAccountFunc = func(username, password string, locked bool) (int, error) {
		return 0, errors.New(expectedErrMsg)
	}

	email := fmt.Sprintf("test_%s@gmail.com", GenerateRandomString(17))
	password := fmt.Sprintf("test_password_%s", GenerateRandomString(17))

	tests := []struct {
		scenario          string
		email             string
		password          string
		res               *proto.CreateAccountResponse
		errCode           codes.Code
		errMsg            string
		ImportAccountFunc func(username, password string, locked bool) (int, error)
	}{
		// scenario: duplicate account
		{
			"account already exists",
			email,
			password,
			nil,
			codes.Unknown,
			expectedErrMsg,
			func(username, password string, locked bool) (int, error) {
				return 0, errors.New(expectedErrMsg)
			},
		},
		// scenario: invalid email params
		{
			"invalid email params",
			"",
			password,
			nil,
			codes.InvalidArgument,
			"invalid input arguments",
			func(username, password string, locked bool) (int, error) {
				return 0, nil
			},
		},
		// scenario: invalid password params
		{
			"invalid password params",
			email,
			"",
			nil,
			codes.InvalidArgument,
			"invalid input arguments",
			func(username, password string, locked bool) (int, error) {
				return 0, nil
			},
		},
		// scenario: valid request
		{
			"valid request",
			email,
			password,
			&proto.CreateAccountResponse{
				Id:    1,
				Error: "",
			},
			codes.Unknown,
			"",
			func(username, password string, locked bool) (int, error) {
				return 1, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.ImportAccountFunc = tt.ImportAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			request := &proto.CreateAccountRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			response, err := client.CreateAccount(ctx, request)

			if response != nil {
				if response.GetError() != tt.res.GetError() {
					t.Error("response: expected", tt.res.GetError(), "received", response.GetError())
				}

				if response.Id != tt.res.Id {
					t.Error("response: expected", tt.res.Id, "received", response.Id)
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
