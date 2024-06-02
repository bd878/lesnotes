package grpc

import (
  "context"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/users/internal/controller"
  "github.com/bd878/gallery/server/users/internal/controller/users"
  "github.com/bd878/gallery/server/users/pkg/model"
)

type Handler struct {
  api.UnimplementedUserServiceServer
  ctrl *users.Controller
}

func New(ctrl *users.Controller) *Handler {
  return &Handler{ctrl: ctrl}
}

func (h *Handler) Auth(ctx context.Context, req *api.AuthUserRequest) (*api.AuthUserResponse, error) {
  if req == nil || req.Token == "" {
    return nil, status.Errorf(codes.InvalidArgument, "nil or empty token")
  }
  u, err := h.ctrl.Get(ctx, &model.User{Token: req.Token})
  if err == controller.ErrTokenInvalid {
    return nil, status.Errorf(codes.InvalidArgument, "wrong token")
  } else if err != nil {
    return nil, status.Errorf(codes.Internal, err.Error())
  }
  return &api.AuthUserResponse{User: model.UserToProto(u)}, nil
}