package grpc

import (
  "context"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usersmodel "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/internal/grpcutil"
)

type Gateway struct {
  userAddr string
}

func New(userAddr string) *Gateway {
  return &Gateway{userAddr}
}

func (g *Gateway) Auth(ctx context.Context, log *logger.Logger, params *model.AuthParams) (*usersmodel.User, error) {
  conn, err := grpcutil.ServiceConnection(ctx, g.userAddr)
  if err != nil {
    log.Error("message", "failed to establish connection with user service")
    return nil, err
  }
  defer conn.Close()
  client := api.NewUserServiceClient(conn)
  resp, err := client.Auth(ctx, &api.AuthUserRequest{Token: params.Token})
  if err != nil {
    log.Error("message", "failed to authenticate on client")
    return nil, err
  }
  return usersmodel.UserFromProto(resp.User), nil
}