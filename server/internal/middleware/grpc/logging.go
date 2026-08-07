package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

type LogReporter struct {
}

func (_ *LogReporter) MsgReceive(req any, info *grpc.UnaryServerInfo, params *MsgReceiveParams) {
	slog.Info("-->", slog.String("method", info.FullMethod), slog.Int64("time", params.Time.UnixMilli()))
}

func (_ *LogReporter) MsgSend(resp any, info *grpc.UnaryServerInfo, params *MsgSendParams) {
	slog.Info("<--", slog.String("method", info.FullMethod), slog.Int64("time", params.Time.UnixMilli()))
	if params.HandlerError != nil {
		slog.Error("handler error", slog.String("error", params.HandlerError.Error()))
	}
}

func LogBuilder() ReporterBuilder {
	return func(ctx context.Context, _ *Meta) (Reporter, context.Context) {
		return &LogReporter{}, ctx
	}
}