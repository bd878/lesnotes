package grpc

import (
  "context"

  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usersmodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Gateway struct {
  userAddr string
}

func New(userAddr string) *Gateway {
  return &Gateway{userAddr}
}

func (g *Gateway) Auth(_ context.Context, _ *logger.Logger, params *model.AuthParams) (*usersmodel.User, error) {
  return &usersmodel.User{
    Id: 1,
    Name: "test",
    Password: "12345",
    Token: "AAAAAA",
    Expires: "never",
  }, nil
}