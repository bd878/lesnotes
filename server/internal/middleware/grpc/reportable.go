package middleware

import (
	"time"
	"context"
	"google.golang.org/grpc"
)

type Meta struct {
}

type MsgReceiveParams struct {
	Time time.Time
}

type MsgSendParams struct {
	Time         time.Time
	HandlerError error
}

type Reporter interface {
	MsgReceive(req any, info *grpc.UnaryServerInfo, params *MsgReceiveParams)
	MsgSend(resp any, info *grpc.UnaryServerInfo, params *MsgSendParams)
}

type ReporterBuilder func(context.Context, *Meta) (Reporter, context.Context)

func (builder ReporterBuilder) Build(ctx context.Context, meta *Meta) (Reporter, context.Context) {
	return builder(ctx, meta)
}