package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/util"
)

func main() {
	serverAddress := "localhost:9999"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}
	defer conn.Close()
	client := proto.NewAuthenticationHandlerServiceApiClient(conn)

	for i := range [10]int{} {
		newAccount := proto.CreateAccountRequest{
			Email:                fmt.Sprintf("yoan_%s@gmail.com", util.GenerateRandomString(20+i)),
			Password:             fmt.Sprintf("%s", util.GenerateRandomString(20+i)),
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}

		if responseMessage, e := client.CreateAccount(context.Background(), &newAccount); e != nil {
			panic(fmt.Sprintf("Was not able to create Record %v", e))
		} else {
			fmt.Println("Record Inserted..")
			fmt.Println(responseMessage)
			fmt.Println("=============================")
		}
	}
}
