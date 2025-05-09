package grpc

import (
	"context"
	"google.golang.org/grpc"

	"github.com/bd878/gallery/server/logger"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
)

type logHandler struct {}

var _ grpcmiddleware.Reporter = (*logHandler)(nil)

func NewBuilder() grpcmiddleware.ReporterBuilder {
	return func(ctx context.Context, _ *grpcmiddleware.Meta) (grpcmiddleware.Reporter, context.Context) {
		return logHandler{}, ctx
	}
}

func (l logHandler) MsgReceive(req any, info *grpc.UnaryServerInfo, params *grpcmiddleware.MsgReceiveParams) {
	logger.Infow("-->", "method", info.FullMethod, "time", params.Time.UnixMilli())
}

func (l logHandler) MsgSend(resp any, info *grpc.UnaryServerInfo, params *grpcmiddleware.MsgSendParams) {
	logger.Infow("<--", "method", info.FullMethod, "time", params.Time.UnixMilli())
	if params.HandlerError != nil {
		logger.Errorln(params.HandlerError.Error())
	}
}
