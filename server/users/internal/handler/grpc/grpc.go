package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type Controller interface {
	GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error)
}

type Handler struct {
	api.UnimplementedUsersServer
	controller Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) Auth(ctx context.Context, req *api.AuthUserRequest) (*api.AuthUserResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil or empty token")
	}
	user, err := h.controller.GetUser(ctx, logger.Default(), &model.GetUserParams{
		Token: req.Token,
	})
	if err == controller.ErrTokenInvalid {
		return nil, status.Errorf(codes.InvalidArgument, "wrong token")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &api.AuthUserResponse{User: model.UserToProto(user)}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {
	logger.Debugw("get user")

	params := &model.GetUserParams{}

	switch key := req.SearchKey.(type) {
	case *api.GetUserRequest_Id:
		logger.Debugw("get user by id", "id", key.Id)
		params.ID = key.Id
	case *api.GetUserRequest_Token:
		logger.Debugw("get user by token", "token", key.Token)
		params.Token = key.Token
	}

	user, err := h.controller.GetUser(ctx, logger.Default(), params)
	switch err {
	case controller.ErrTokenInvalid:
		return nil, status.Errorf(codes.InvalidArgument, "wrong token")
	case nil:
		logger.Debugw("ger user ok", "id", user.ID)
		return model.UserToProto(user), nil
	default:
		return nil, status.Errorf(codes.Internal, err.Error())
	}
}