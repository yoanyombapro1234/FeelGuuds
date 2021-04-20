package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/util"
)

func TestAuthenticationHandlerServiceClient_GetAccount(t *testing.T) {
	email := fmt.Sprintf("test_%s@gmail.com", util.GenerateRandomString(17))
	password := fmt.Sprintf("test_%s", util.GenerateRandomString(17))

	tests := []struct {
		name string
		email   string
		password string
		err    error
		accountExists bool
	}{
		{
			"get non-existent account scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 404 from http://localhost:8000/accounts/1000. Errors in account: NOT_FOUND"),
			false,

		},
		{
			"get account that already exists scenario",
			email,
			password,
			nil,
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
			var account *proto.Account
			if !tt.accountExists {
				account, err = c.GetAccount(context.Background(), NON_EXISTENT_ACCOUNT_ID)
			} else {
				account, err = c.GetAccount(context.Background(), id)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}

			if err == nil && account != nil && account.Username != tt.email {
				t.Error("error: expected", tt.email, "received", account.Username)
			}
		})
	}
}
