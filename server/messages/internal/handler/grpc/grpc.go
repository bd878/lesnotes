package grpc

import (
  "context"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Controller interface {
  SaveMessage(ctx context.Context, msg *model.Message) (*model.Message, error)
  ReadUserMessages(ctx context.Context, userId usermodel.UserId, limit int32, offset int32, ascending bool) (*model.MessagesList, error)
  GetServers(ctx context.Context) ([]*api.Server, error)
}

type Handler struct {
  api.UnimplementedMessagesServer
  ctrl Controller
}

func New(ctrl Controller) *Handler {
  h := &Handler{ctrl: ctrl}

  return h
}

func (h *Handler) SaveMessage(ctx context.Context, req *api.SaveMessageRequest) (
  *api.SaveMessageResponse,
  error,
) {
  msg, err := h.ctrl.SaveMessage(ctx, model.MessageFromProto(req.Message))
  if err != nil {
    return &api.SaveMessageResponse{Message: req.Message}, err
  }
  return &api.SaveMessageResponse{Message: model.MessageToProto(msg)}, nil
}

func (h *Handler) ReadUserMessages(ctx context.Context, req *api.ReadUserMessagesRequest) (
  *api.ReadUserMessagesResponse,
  error,
) {
  var res *model.MessagesList
  var err error

  res, err = h.ctrl.ReadUserMessages(
    ctx,
    usermodel.UserId(req.UserId),
    req.Limit,
    req.Offset,
    req.Asc,
  )
  if err != nil {
    return nil, err
  }

  return &api.ReadUserMessagesResponse{
    Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
    IsLastPage: res.IsLastPage,
  }, nil
}

func (h *Handler) GetServers(ctx context.Context, req *api.GetServersRequest) (
  *api.GetServersResponse,
  error,
) {
  srvs, err := h.ctrl.GetServers(ctx)
  if err != nil {
    return &api.GetServersResponse{Servers: nil}, err
  }
  return &api.GetServersResponse{
    Servers: srvs,
  }, nil
}
