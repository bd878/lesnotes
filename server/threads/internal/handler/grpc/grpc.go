package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Controller interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*threads.Thread, isLastPage bool, err error)
	ReadThread(ctx context.Context, id, userID int64) (thread *threads.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error)
	UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error)
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	PublishThread(ctx context.Context, id, userID int64) (err error)
	PrivateThread(ctx context.Context, id, userID int64) (err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedThreadsServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
}

func (h *Handler) Read(ctx context.Context, req *api.ReadRequest) (resp *api.Thread, err error) {
	thread, err := h.controller.ReadThread(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = threads.ThreadToProto(thread)

	return resp, nil
}

func (h *Handler) List(ctx context.Context, req *api.ListRequest) (resp *api.ListResponse, err error) {
	list, isLastPage, err := h.controller.ListThreads(ctx, req.UserId, req.ParentId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	resp = &api.ListResponse{
		List:       threads.MapThreadsToProto(threads.ThreadToProto, list),
		IsLastPage: isLastPage,
		Count:      int32(len(list)),
	}

	return
}

func (h *Handler) Resolve(ctx context.Context, req *api.ResolveRequest) (resp *api.ResolveResponse, err error) {
	ids, err := h.controller.ResolveThread(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.ResolveResponse{
		Path: ids,
	}

	return
}

func (h *Handler) Create(ctx context.Context, req *api.CreateRequest) (resp *api.CreateResponse, err error) {
	// TODO: validate that parent thread exists

	err = h.controller.CreateThread(ctx, req.Id, req.UserId, req.ParentId, req.NextId, req.PrevId, req.Name, req.Description, req.Private)
	if err != nil {
		return
	}

	resp = &api.CreateResponse{}

	return
}

func (h *Handler) Update(ctx context.Context, req *api.UpdateRequest) (resp *api.UpdateResponse, err error) {
	err = h.controller.UpdateThread(ctx, req.Id, req.UserId, req.Name, req.Description)
	if err != nil {
		return
	}

	resp = &api.UpdateResponse{}

	return
}

func (h *Handler) Reorder(ctx context.Context, req *api.ReorderRequest) (resp *api.ReorderResponse, err error) {
	err = h.controller.ReorderThread(ctx, req.Id, req.UserId, req.ParentId, req.NextId, req.PrevId)
	if err != nil {
		return
	}

	resp = &api.ReorderResponse{}

	return
}

func (h *Handler) Delete(ctx context.Context, req *api.DeleteRequest) (resp *api.DeleteResponse, err error) {
	err = h.controller.DeleteThread(ctx, req.Id, req.UserId)
	if err != nil {
		return
	}

	resp = &api.DeleteResponse{}

	return
}

func (h *Handler) Publish(ctx context.Context, req *api.PublishRequest) (resp *api.PublishResponse, err error) {
	err = h.controller.PublishThread(ctx, req.Id, req.UserId)
	if err != nil {
		return
	}

	resp = &api.PublishResponse{}

	return
}

func (h *Handler) Private(ctx context.Context, req *api.PrivateRequest) (resp *api.PrivateResponse, err error) {
	err = h.controller.PrivateThread(ctx, req.Id, req.UserId)
	if err != nil {
		return
	}

	resp = &api.PrivateResponse{}

	return
}

func (h *Handler) Count(ctx context.Context, req *api.CountRequest) (resp *api.CountResponse, err error) {
	total, err := h.controller.CountThreads(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.CountResponse{
		Total: total,
	}

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
