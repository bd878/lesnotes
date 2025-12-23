package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/search/pkg/model"
)

type Controller interface {
	SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*model.Message, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedSearchServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
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
