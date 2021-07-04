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
)

func Test_delete_account(t *testing.T) {
	// TODO : ensure proper metrics are being emitted in each unit test
	expectedErrMsg := "retry limit reached"

	id := "1"

	tests := []struct {
		scenario           string
		id                 string
		res                *proto.DeleteAccountResponse
		errCode            codes.Code
		errMsg             string
		ArchiveAccountFunc func(id string) error
	}{
		// scenario: duplicate account
		{
			"account already exists",
			id,
			nil,
			codes.Unknown,
			expectedErrMsg,
			func(id string) error {
				return errors.New("account not found")
			},
		},
		// scenario: valid request
		{
			"valid request",
			id,
			nil,
			codes.Unknown,
			"",
			func(id string) error {
				return nil
			},
		},
		// scenario: invalid request
		{
			"valid request",
			"0",
			nil,
			codes.InvalidArgument,
			"invalid input argument",
			func(id string) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()
			ThirdPartyMockService.ArchiveAccountFunc = tt.ArchiveAccountFunc
			conn := MockGRPCService(ctx, &ThirdPartyMockService)
			defer conn.Close()

			client := proto.NewAuthenticationHandlerServiceApiClient(conn)

			accountId, _ := strconv.Atoi(tt.id)
			request := &proto.DeleteAccountRequest{
				Id: uint32(accountId),
			}

			response, err := client.DeleteAccount(ctx, request)

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
