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

func TestAuthenticationHandlerServiceClient_UpdateAccount(t *testing.T) {
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_%s", util.GenerateRandomString(17))
	updatedEmail := email + util.GenerateRandomString(11)

	tests := []struct {
		name string
		email   string
		err    error
		accountExist bool
	}{
		{
			"update account scenario that doesnt exist",
			updatedEmail,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 404 from http://localhost:8000/accounts/10000. Errors in account: NOT_FOUND"),
			false,
		},
		{
			"update account that already exists scenario",
			updatedEmail,
			nil,
			true,
		},
		{
			"update account with invalid parameters scenario (missing email)",
			"",
			errors.New("grpc: InvalidArgument, invalid input arguments"),
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
			if !tt.accountExist {
				err = c.UpdateAccount(context.Background(), NON_EXISTENT_ACCOUNT_ID, tt.email)
			} else {
				err = c.UpdateAccount(context.Background(), id, tt.email)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}
}
