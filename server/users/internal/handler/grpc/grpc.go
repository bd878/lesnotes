package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Controller interface {
	GetUser(ctx context.Context, id int32) (*model.User, error)
	FindUser(ctx context.Context, params *model.FindUserParams) (*model.User, error)
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
	params := &model.FindUserParams{}

	switch key := req.SearchKey.(type) {
	case *api.FindUserRequest_Name:
		params.Name = key.Name
	case *api.FindUserRequest_Token:
		params.Token = key.Token
	}

	user, err := h.controller.FindUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return model.UserToProto(user), nil
}