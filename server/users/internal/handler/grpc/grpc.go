package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Controller interface {
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, newLogin string, metadata []byte) (err error)
	FindUser(ctx context.Context, login string) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) (err error)
	CreateUser(ctx context.Context, id int64, login, password string, metadata []byte) (*model.User, error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedUsersServer
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {
	user, err := h.controller.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return model.UserToProto(user), nil
}

func (h *Handler) FindUser(ctx context.Context, req *api.FindUserRequest) (*api.User, error) {
	user, err := h.controller.FindUser(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	return model.UserToProto(user), nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (resp *api.UpdateUserResponse, err error) {
	err = h.controller.UpdateUser(ctx, req.Id, req.Login, req.Metadata)
	if err != nil {
		return
	}

	resp = &api.UpdateUserResponse{}

	return
}

func (h *Handler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (resp *api.CreateUserResponse, err error) {
	_, err = h.controller.CreateUser(ctx, req.Id, req.Login, req.Password, req.Metadata)
	if err != nil {
		return
	}

	resp = &api.CreateUserResponse{}

	return
}

func (h *Handler) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (resp *api.DeleteUserResponse, err error) {
	err = h.controller.DeleteUser(ctx, req.Id)
	if err != nil {
		return
	}

	resp = &api.DeleteUserResponse{}

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
