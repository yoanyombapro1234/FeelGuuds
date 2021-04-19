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

// TODO: ensure the authentication client adheres to an interface which we can directly mock and pass in
// a version so we can better test various use cases without having to spin up the service dependencies - https://dev.to/jonfriesen/mocking-dependencies-in-go-1h4d
// also add unit and integration testing scenarios at the function level
func Test_create_account(t *testing.T) {
	expectedErrMsg := "retry limit reached"
	ThirdPartyMockService.ImportAccountFunc = func(username, password string, locked bool) (int, error) {
		return 0, errors.New(expectedErrMsg)
	}

	email := fmt.Sprintf("test_%s@gmail.com", GenerateRandomString(17))
	password := fmt.Sprintf("test_password_%s", GenerateRandomString(17))

	tests := []struct {
		scenario string
		email    string
		password string
		res      *proto.CreateAccountResponse
		errCode  codes.Code
		errMsg   string
	}{
		// scenario: valid fields
		{
			"account already exists",
			email,
			password,
			nil,
			codes.Unknown,
			expectedErrMsg,
		},
	}

	ctx := context.Background()
	conn := MockGRPCService(ctx, &ThirdPartyMockService)
	defer conn.Close()

	client := proto.NewAuthenticationHandlerServiceApiClient(conn)

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			request := &proto.CreateAccountRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			response, err := client.CreateAccount(ctx, request)

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
