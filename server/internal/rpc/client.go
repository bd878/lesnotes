package rpc

import (
	"time"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	return grpc.NewClient(target,
		append([]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time: 5*time.Minute,
				Timeout: 10*time.Second,
				PermitWithoutStream: true,
			}),
		}, opts...)...,
	)
}