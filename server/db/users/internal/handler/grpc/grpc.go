package grpc

import (
	"time"
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/db/users/pkg/machine"
)

type Controller interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	GetUser(ctx context.Context, id int64) (*users.User, error)
	FindUser(ctx context.Context, login string) (*users.User, error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	users.UnimplementedUsersServer
	api.UnimplementedDistributedServer
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller: controller}
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

func (h *Handler) GetUser(ctx context.Context, req *users.GetUserRequest) (*users.User, error) {
	user, err := h.controller.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) FindUser(ctx context.Context, req *users.FindUserRequest) (*users.User, error) {
	user, err := h.controller.FindUser(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	return user, nil
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
