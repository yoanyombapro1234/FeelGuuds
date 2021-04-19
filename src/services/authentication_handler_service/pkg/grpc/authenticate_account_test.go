package grpc

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
	"errors"
)

func Test_authenticate_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	email := fmt.Sprintf("test_%s@gmail.com", GenerateRandomString(17))
	password := fmt.Sprintf("test_password_%s", GenerateRandomString(17))

	tests := []struct {
		scenario string
		email    string
		password string
		res      *proto.CreateAccountResponse
		errCode  codes.Code
		errMsg   string
		LoginAccount func(username, password string) (string, error)
	}{
		// scenario: invalid input arguments - password
		{
			"invalid input arguments - password",
			email,
			"",
			nil,
			codes.InvalidArgument,
			service_errors.ErrInvalidInputArguments.Error(),
			func(username, password string) (string, error) {
				return "", service_errors.ErrInvalidInputArguments
			},
		},
		// scenario: invalid input params - email
		{
			"invalid email params",
			"",
			password,
			nil,
			codes.InvalidArgument,
			service_errors.ErrInvalidInputArguments.Error(),
			func(username, password string) (string, error) {
				return "", service_errors.ErrInvalidInputArguments
			},
		},
		// scenario: valid request
		{
			"valid request",
			email,
			password,
			nil,
			codes.Unknown,
			"",
			func(username, password string) (string, error) {
				return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", nil
			},
		},
		// scenario: invalid input params - email
		{
			"bad credentials",
			email,
			password,
			nil,
			codes.Unknown,
			"retry limit reached (1/1): field: credentials, message: failed",
			func(username, password string) (string, error) {
				return "", errors.New("field: credentials, message: failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.LoginAccountFunc = tt.LoginAccount
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			request := &proto.AuthenticateAccountRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			response, err := client.AuthenticateAccount(ctx, request)

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
					if !strings.Contains(tt.errMsg, er.Message()) {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
