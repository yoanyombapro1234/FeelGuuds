package client

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestAuthenticationHandlerServiceClient_CreateAccount(t *testing.T) {
	email := "test_yoan@gmail.com" // fmt.Sprintf("test_yoan@gmail.com", util.GenerateRandomString(17))
	password := "test_yoan" //fmt.Sprintf("test_yoan", util.GenerateRandomString(17))

	tests := []struct {
		name     string
		email    string
		password string
		err      error
	}{
		{
			"create account scenario",
			email,
			password,
			nil,
		},
		{
			"create account that already exists scenario",
			email,
			password,
			errors.New("grpc: Unknown, retry limit reached (1/1): received 422 from http://localhost:8000/accounts/import. Errors in username: TAKEN"),
		},
		{
			"create account with invalid parameters scenario (missing email)",
			"",
			password,
			errors.New("grpc: InvalidArgument, invalid input arguments"),
		},
		{
			"create account with invalid parameters scenario (missing password)",
			email,
			"",
			errors.New("grpc: InvalidArgument, invalid input arguments"),
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
			c := NewClient(conn, 5*time.Second)
			_, err := createAccountTestHelper(c, ctx, tt.email, tt.password)
			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}
}
