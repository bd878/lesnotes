package grpc

import (
  "context"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type Controller interface {
  SaveMessage(ctx context.Context, msg *model.Message) (model.MessageId, error)
  ReadUserMessages(ctx context.Context, userId usermodel.UserId) ([]*model.Message, error)
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
  msgId, err := h.ctrl.SaveMessage(ctx, model.MessageFromProto(req.Message))
  if err != nil {
    return &api.SaveMessageResponse{Id: uint32(0)}, err
  }
  return &api.SaveMessageResponse{Id: uint32(msgId)}, nil
}

func (h *Handler) ReadUserMessages(ctx context.Context, req *api.ReadUserMessagesRequest) (
  *api.ReadUserMessagesResponse,
  error,
) {
  msgs, err := h.ctrl.ReadUserMessages(ctx, usermodel.UserId(req.UserId))
  if err != nil {
    return nil, err
  }
  return &api.ReadUserMessagesResponse{
    Messages: model.MapMessagesToProto(model.MessageToProto, msgs),
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
