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

func TestAuthenticationHandlerServiceClient_DeleteAccount(t *testing.T) {
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_%s", util.GenerateRandomString(17))

	tests := []struct {
		name string
		email   string
		password string
		err    error
		shouldCreateAccount  bool
		shouldInputArgBeMisconfigured bool
	}{
		{
			"delete non-existent account scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 404 from http://localhost:8000/accounts/1000. Errors in account: NOT_FOUND"),
			false,
			false,
		},
		{
			"delete account that already exists scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 422 from http://localhost:8000/accounts/import. Errors in username: TAKEN"),
			true,
			false,
		},
		{
			"attempt to delete account with invalid parameters scenario (missing id)",
			"",
			password,
			errors.New("grpc: InvalidArgument, invalid input arguments"),
			false,
			true,
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
			var id uint32 = 1000
			c :=  NewClient(conn, 5*time.Second)
			if tt.shouldCreateAccount {
				id, err = createAccountTestHelper(c, ctx, tt.email, tt.password)
				if err != nil {
					t.Error("error: expected", tt.err, "received", err)
				}
			}

			if tt.shouldInputArgBeMisconfigured {
				err = c.ArchiveAccount(context.Background(), 0)
			} else {
				err = c.ArchiveAccount(context.Background(), id)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}
}
