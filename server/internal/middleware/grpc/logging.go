package middleware

import (
	"context"
	"github.com/bd878/gallery/server/logger"
)

type LogReporter struct {
}

func (_ *LogReporter) MsgReceive(_ any, _ *MsgReceiveParams) {
	logger.Infoln("request start")
}

func (_ *LogReporter) MsgSend(_ any, _ *MsgSendParams) {
	logger.Infoln("request end")
}

func LogBuilder() ReporterBuilder {
	return func(ctx context.Context, _ *Meta) (
		Reporter, context.Context,
	) {
		return &LogReporter{}, ctx
	}
}
