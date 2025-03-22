package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(_ context.Context, addr string) (*grpc.ClientConn, error) {
	// todo: service registry
	return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}