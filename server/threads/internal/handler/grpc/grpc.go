package grpc

import (
	"time"
	"errors"
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/threads/internal/machine"
	"github.com/bd878/gallery/server/threads/internal/controller"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Controller interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*threads.Thread, isLastPage bool, err error)
	ListMessages(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool, privateMessage *bool) (ids []*threads.Thread, isLastPage bool, err error)
	ReadThread(ctx context.Context, id, userID int64, name string) (thread *threads.Thread, err error)
	ReadParent(ctx context.Context, id, userID int64) (thread *threads.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (path []*api.PathStep, err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	CountMessages(ctx context.Context, id, userID int64, privateMessage *bool) (total int32, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedThreadsServer
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

func (h *Handler) Read(ctx context.Context, req *api.ReadRequest) (resp *api.Thread, err error) {
	thread, err := h.controller.ReadThread(ctx, req.Id, req.UserId, req.Name)
	if err != nil {
		return nil, err
	}

	resp = threads.ThreadToProto(thread)

	return
}

func (h *Handler) List(ctx context.Context, req *api.ListRequest) (resp *api.ListResponse, err error) {
	list, isLastPage, err := h.controller.ListThreads(ctx, req.UserId, req.ParentId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	// TODO: do not list child threads of private parent thread.
	// But it should list all threads (public and private) of a public thread.

	resp = &api.ListResponse{
		List:       threads.MapThreadsToProto(threads.ThreadToProto, list),
		IsLastPage: isLastPage,
		Count:      int32(len(list)),
	}

	return
}

func (h *Handler) ListMessages(ctx context.Context, req *api.ListMessagesRequest) (resp *api.ListMessagesResponse, err error) {
	list, isLastPage, err := h.controller.ListMessages(ctx, req.UserId, req.ParentId, req.Limit, req.Offset, req.Asc, req.Private)
	if err != nil {
		return nil, err
	}

	resp = &api.ListMessagesResponse{
		List:       threads.MapThreadsToProto(threads.ThreadToProto, list),
		IsLastPage: isLastPage,
		Count:      int32(len(list)),
	}

	return
}

func (h *Handler) Resolve(ctx context.Context, req *api.ResolveRequest) (resp *api.ResolveResponse, err error) {
	path, err := h.controller.ResolveThread(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.ResolveResponse{
		Path: path,
	}

	return
}

func (h *Handler) ReadParent(ctx context.Context, req *api.ReadParentRequest) (resp *api.ReadParentResponse, err error) {
	parent, err := h.controller.ReadParent(ctx, req.Id, req.UserId)
	if err != nil {
		if errors.Is(err, controller.ErrParentIsRoot) {
			return &api.ReadParentResponse{IsRoot: true}, nil
		} else {
			return nil, err
		}
	}

	resp = &api.ReadParentResponse{Parent: threads.ThreadToProto(parent), IsRoot: false}

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

func (h *Handler) CountMessages(ctx context.Context, req *api.CountMessagesRequest) (resp *api.CountMessagesResponse, err error) {
	total, err := h.controller.CountMessages(ctx, req.Id, req.UserId, req.PrivateMessage)
	if err != nil {
		return nil, err
	}

	resp = &api.CountMessagesResponse{
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
