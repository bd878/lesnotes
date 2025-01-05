package middleware

import (
  "context"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/internal/middleware"
)

type LogReporter struct {
}

func (_ *LogReporter) MsgReceive(_ any, _ *middleware.MsgReceiveParams) {
  logger.Infoln("request start")
}

func (_ *LogReporter) MsgSend(_ any, _ *middleware.MsgSendParams) {
  logger.Infoln("request end")
}

func LogBuilder() middleware.ReporterBuilder {
  return func(ctx context.Context, _ *middleware.Meta) (
    middleware.Reporter, context.Context,
  ) {
    return &LogReporter{}, ctx
  }
}
