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
				Timeout: 5*time.Second,
				PermitWithoutStream: false,
			}),
			grpc.WithIdleTimeout(2 * time.Minute),
		}, opts...)...,
	)
}