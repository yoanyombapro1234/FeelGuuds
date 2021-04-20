package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	client2 "github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/client"
)

func main() {
	log.Println("Client running ...")

	conn, err := grpc.Dial(":9999", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := client2.NewClient(conn, time.Duration(20))
	response, err := client.GetAccount(context.Background(), 1)
	log.Println(response)
	log.Println(err)
}
