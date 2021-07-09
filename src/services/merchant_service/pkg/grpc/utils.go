package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func (s *Server) ConfigureOutgoingRpcRequest(ctx context.Context) (context.Context, context.CancelFunc) {
	md := metadata.Pairs("origin_svc", s.config.ServiceName)
	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx, cancel := s.setCtxRequestTimeout(ctx)
	return ctx, cancel
}
