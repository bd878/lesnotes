package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
)

type Controller interface {
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name string, private bool) (err error)
	UpdateThread(ctx context.Context, id, userID int64, name string, private int32) (err error)
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	PublishThread(ctx context.Context, id, userID int64) (err error)
	PrivateThread(ctx context.Context, id, userID int64) (err error)
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

func (h *Handler) Resolve(ctx context.Context, req *api.ResolveRequest) (resp *api.ResolveResponse, err error) {
	ids, err := h.controller.ResolveThread(ctx, req.UserId, req.UserId)
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

	err = h.controller.CreateThread(ctx, req.Id, req.UserId, req.ParentId, req.NextId, req.PrevId, req.Name, req.Private)
	if err != nil {
		return
	}

	resp = &api.CreateResponse{}

	return
}

func (h *Handler) Update(ctx context.Context, req *api.UpdateRequest) (resp *api.UpdateResponse, err error) {
	err = h.controller.UpdateThread(ctx, req.Id, req.UserId, req.Name, req.Private)
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
