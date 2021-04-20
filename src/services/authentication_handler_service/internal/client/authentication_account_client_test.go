package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/util"
)

func TestAuthenticationHandlerServiceClient_AuthenticateAccount(t *testing.T) {
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_%s", util.GenerateRandomString(17))

	tests := []struct {
		name string
		email   string
		password string
		err    error
		shouldCreateAccount  bool
		shouldCheckToken bool
	}{
		{
			"authenticate non-existent account scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 422 from http://localhost:8000/session. Errors in credentials: FAILED"),
			false,
			false,
		},
		{
			"authenticate account that already exists scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 422 from http://localhost:8000/accounts/import. Errors in username: TAKEN"),
			true,
			true,
		},
		{
			"attempt to authenticate account with invalid parameters scenario (missing email)",
			"",
			password,
			errors.New("grpc: InvalidArgument, invalid input arguments"),
			false,
			false,
		},
		{
			"attempt to authenticate account with invalid parameters scenario (missing password)",
			email,
			"",
			errors.New("grpc: InvalidArgument, invalid input arguments"),
			false,
			false,
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c :=  NewClient(conn, 5*time.Second)
			if tt.shouldCreateAccount {
				_, err := createAccountTestHelper(c, ctx, tt.email, tt.password)
				if err != nil {
					t.Error("error: expected", tt.err, "received", err)
				}
			}

			token, err := c.AuthenticateAccount(context.Background(), tt.email, tt.password)

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}

			if tt.shouldCheckToken {
				if token == "" {
					t.Error("error: expected non empty jwt token")
				}
			}
		})
	}
}
