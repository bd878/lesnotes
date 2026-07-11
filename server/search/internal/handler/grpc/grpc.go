package grpc

import (
	"time"
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/search/pkg/model"
	"github.com/bd878/gallery/server/search/internal/machine"
)

type Controller interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*model.Message, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedSearchServer
	api.UnimplementedDistributedServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
}

func (h *Handler) Apply(ctx context.Context, req *api.Command) (resp *api.CommandResponse, err error) {
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	err = h.controller.Apply(ctx, machine.RequestType(req.ReqType), req.Cmd, duration)

	resp = &api.CommandResponse{}

	return
}

func (h *Handler) SearchMessages(ctx context.Context, req *api.SearchMessagesRequest) (resp *api.SearchMessagesResponse, err error) {
	var (
		public int
		threadID int64
	)

	if req.ThreadId == nil {
		threadID = 0
	} else {
		threadID = *req.ThreadId
	}

	if req.Public == nil {
		public = -1
	} else {
		public = int(*req.Public)
	}

	list, err := h.controller.SearchMessages(ctx, req.UserId, req.Substr, threadID, public)
	if err != nil {
		return nil, err
	}

	logger.Debugw("found messages", "count", len(list))

	resp = &api.SearchMessagesResponse{
		List:   model.MapMessagesToProto(model.MessageToProto, list),
		Count:  int32(len(list)),
	}

	return
}

func (h *Handler) SearchFiles(ctx context.Context, req *api.SearchFilesRequest) (resp *api.SearchFilesResponse, err error) {
	// not implemented
	return
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (resp *api.GetServersResponse, err error) {
	servers, err := h.controller.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	resp = &api.GetServersResponse{
		Servers: servers,
	}

	return
}
