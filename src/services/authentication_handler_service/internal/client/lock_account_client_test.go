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

func TestAuthenticationHandlerServiceClient_LockAccount(t *testing.T) {
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_%s", util.GenerateRandomString(17))

	tests := []struct {
		name string
		email   string
		password string
		err    error
		accountExists bool
		shouldAccountBeLockPrematurely bool
	}{
		{
			"lock non-existent account scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 404 from http://localhost:8000/accounts/1000. Errors in account: NOT_FOUND"),
			false,
			false,

		},
		{
			"lock account that already exists and isnt locked scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 422 from http://localhost:8000/accounts/import. Errors in username: TAKEN"),
			true,
			false,
		},
		{
			"lock account that already exists and is locked scenario",
			"",
			password,
			nil,
			true,
			true,
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c :=  NewClient(conn, 5*time.Second)

	id, err := createAccountTestHelper(c, ctx, email, password)
	if err != nil {
		t.Error("error: did not expect to receive an error")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.accountExists {
				err = c.LockAccount(context.Background(), NON_EXISTENT_ACCOUNT_ID)
			} else {
				if tt.shouldAccountBeLockPrematurely {
					_ = c.LockAccount(context.Background(), id)
				}

				err = c.LockAccount(context.Background(), id)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}
}
