package middleware

import (
	"time"
	"context"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(builder ReporterBuilder) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		reporter, newCtx := builder.Build(ctx, nil)
		reporter.MsgReceive(req, info, &MsgReceiveParams{
			Time: time.Now(),
		})

		resp, err := handler(newCtx, req)

		reporter.MsgSend(resp, info, &MsgSendParams{
			Time:  time.Now(),
			HandlerError: err,
		})
		return resp, err
	}
}