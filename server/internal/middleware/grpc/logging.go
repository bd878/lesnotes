package middleware

import (
	"context"
	"github.com/bd878/gallery/server/internal/logger"
	"google.golang.org/grpc"
)

type LogReporter struct {
}

func (_ *LogReporter) MsgReceive(req any, info *grpc.UnaryServerInfo, params *MsgReceiveParams) {
	logger.Infow("-->", "method", info.FullMethod, "time", params.Time.UnixMilli())
}

func (_ *LogReporter) MsgSend(resp any, info *grpc.UnaryServerInfo, params *MsgSendParams) {
	logger.Infow("<--", "method", info.FullMethod, "time", params.Time.UnixMilli())
	if params.HandlerError != nil {
		logger.Errorln(params.HandlerError.Error())
	}
}

func LogBuilder() ReporterBuilder {
	return func(ctx context.Context, _ *Meta) (Reporter, context.Context) {
		return &LogReporter{}, ctx
	}
}
