package grpc

import (
  "context"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/user/pkg/model"
  "github.com/bd878/gallery/server/internal/grpcutil"
)

type Gateway struct {
  userAddr string
}

func New(userAddr string) *Gateway {
  return &Gateway{userAddr}
}

func (g *Gateway) Auth(ctx context.Context, token string) (*model.User, error) {
  conn, err := grpcutil.ServiceConnection(ctx, g.userAddr)
  if err != nil {
    return nil, err
  }
  defer conn.Close()
  client := api.NewUserServiceClient(conn)
  resp, err := client.Auth(ctx, &api.AuthUserRequest{Token: token})
  if err != nil {
    return nil, err
  }
  return model.UserFromProto(resp.User), nil
}