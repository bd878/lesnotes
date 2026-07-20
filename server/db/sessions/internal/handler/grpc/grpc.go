package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/sessions"
)

type Controller interface {
	GetSession(ctx context.Context, token string) (*sessions.Session, error)
	ListUserSessions(ctx context.Context, userID int64) ([]*sessions.Session, error)
	CreateSession(ctx context.Context, userID int64) (*sessions.Session, error)
	RemoveSession(ctx context.Context, token string) error
	RemoveUserSessions(ctx context.Context, userID int64) error
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	sessions.UnimplementedSessionsServer
	api.UnimplementedDistributedServer
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) List(ctx context.Context, req *sessions.ListUserSessionsRequest) (resp *sessions.ListUserSessionsResponse, err error) {
	list, err := h.controller.ListUserSessions(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &sessions.ListUserSessionsResponse{
		Sessions: list,
	}

	return
}

func (h *Handler) Get(ctx context.Context, req *sessions.GetSessionRequest) (resp *sessions.Session, err error) {
	resp, err = h.controller.GetSession(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return
}

func (h *Handler) Create(ctx context.Context, req *sessions.CreateSessionRequest) (resp *sessions.Session, err error) {
	resp, err = h.controller.CreateSession(ctx, req.UserId)
	if err != nil {
		return
	}

	return
}

func (h *Handler) Remove(ctx context.Context, req *sessions.RemoveSessionRequest) (resp *sessions.RemoveSessionResponse, err error) {
	err = h.controller.RemoveSession(ctx, req.Token)
	if err != nil {
		return
	}

	resp = &sessions.RemoveSessionResponse{}

	return
}

func (h *Handler) RemoveAll(ctx context.Context, req *sessions.RemoveAllSessionsRequest) (resp *sessions.RemoveAllSessionsResponse, err error) {
	err = h.controller.RemoveUserSessions(ctx, req.UserId)
	if err != nil {
		return
	}

	resp = &sessions.RemoveAllSessionsResponse{}
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
