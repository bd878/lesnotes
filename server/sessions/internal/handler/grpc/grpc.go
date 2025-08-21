package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

type Controller interface {
	GetSession(ctx context.Context, token string) (*model.Session, error)
	ListUserSessions(ctx context.Context, userID int64) ([]*model.Session, error)
	CreateSession(ctx context.Context, userID int64) (*model.Session, error)
	RemoveSession(ctx context.Context, token string) error
	RemoveUserSessions(ctx context.Context, userID int64) error
}

type Handler struct {
	api.UnimplementedSessionsServer
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) List(ctx context.Context, req *api.ListUserSessionsRequest) (resp *api.ListUserSessionsResponse, err error) {
	sessions, err := h.controller.ListUserSessions(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &api.ListUserSessionsResponse{
		Sessions: model.MapSessionsToProto(model.SessionToProto, sessions),
	}

	return
}

func (h *Handler) Get(ctx context.Context, req *api.GetSessionRequest) (resp *api.Session, err error) {
	session, err := h.controller.GetSession(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	resp = model.SessionToProto(session)

	return
}

func (h *Handler) Create(ctx context.Context, req *api.CreateSessionRequest) (resp *api.Session, err error) {
	var session *model.Session

	session, err = h.controller.CreateSession(ctx, req.UserId)
	if err != nil {
		return
	}

	resp = model.SessionToProto(session)

	return
}

func (h *Handler) Remove(ctx context.Context, req *api.RemoveSessionRequest) (resp *api.RemoveSessionResponse, err error) {
	err = h.controller.RemoveSession(ctx, req.Token)
	if err != nil {
		return
	}

	resp = &api.RemoveSessionResponse{}

	return
}

func (h *Handler) RemoveAll(ctx context.Context, req *api.RemoveAllSessionsRequest) (resp *api.RemoveAllSessionsResponse, err error) {
	err = h.controller.RemoveUserSessions(ctx, req.UserId)
	if err != nil {
		return
	}

	resp = &api.RemoveAllSessionsResponse{}
	return
}
