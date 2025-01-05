package middleware

import (
  "time"
  "context"
  "google.golang.org/grpc"
  "github.com/bd878/gallery/server/messages/internal/middleware"
)

func UnaryServerInterceptor(builder middleware.ReporterBuilder) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
    reporter, newCtx := builder.Build(ctx, nil)
    reporter.MsgReceive(req, &middleware.MsgReceiveParams{
      Time: time.Now(),
    })

    resp, err := handler(newCtx, req)

    reporter.MsgSend(resp, &middleware.MsgSendParams{
      Time:  time.Now(),
      HandlerError: err,
    })
    return resp, err
  }
}