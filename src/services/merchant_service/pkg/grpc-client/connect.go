package grpc_client

import (
	"context"

	grpc_utils "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc-utils"
	tlscert "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/tlsCert"
	"google.golang.org/grpc"
)

func ConnectToClient(addr string) (*grpc.ClientConn, error) {
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithInsecure()
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithStreamInterceptors(grpc_utils.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(grpc_utils.GetDefaultUnaryClientInterceptors())
	cc, err := clientBuilder.GetConn(addr)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func ConnectToClientWithTls(addr string) (*grpc.ClientConn, error) {
	clientBuilder := GrpcConnBuilder{}
	clientBuilder.WithContext(context.Background())
	clientBuilder.WithClientTransportCredentials(false, tlscert.CertPool)
	clientBuilder.WithStreamInterceptors(grpc_utils.GetDefaultStreamClientInterceptors())
	clientBuilder.WithUnaryInterceptors(grpc_utils.GetDefaultUnaryClientInterceptors())
	cc, err := clientBuilder.GetConn(addr)
	if err != nil {
		return nil, err
	}

	return cc, nil
}
