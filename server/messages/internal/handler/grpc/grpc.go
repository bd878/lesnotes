package grpc

import (
  "time"
  "context"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Controller interface {
  SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) error
  UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error
  DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error
  ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (*model.ReadUserMessagesResult, error)
  GetServers(ctx context.Context, log *logger.Logger) ([]*api.Server, error)
}

type Handler struct {
  api.UnimplementedMessagesServer
  controller Controller
}

func New(ctrl Controller) *Handler {
  handler := &Handler{controller: ctrl}

  return handler
}

func (h *Handler) SaveMessage(ctx context.Context, req *api.SaveMessageRequest) (
  *api.SaveMessageResponse, error,
) {
  req.Message.CreateUtcNano = time.Now().UnixNano()
  req.Message.UpdateUtcNano = time.Now().UnixNano()
  req.Message.Id = utils.RandomID()

  err := h.controller.SaveMessage(ctx, logger.Default(), &model.SaveMessageParams{
    Message: model.MessageFromProto(req.Message),
  }) 
  if err != nil {
    logger.Error("message", "failed to save message")
    return nil, err
  }

  return &api.SaveMessageResponse{
    Id: req.Message.Id,
    CreateUtcNano: req.Message.CreateUtcNano,
    UpdateUtcNano: req.Message.UpdateUtcNano,
  }, nil
}

func (h *Handler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (
  *api.UpdateMessageResponse, error,
) {
  updateUTCNano := time.Now().UnixNano()

  err := h.controller.UpdateMessage(ctx, logger.Default(), &model.UpdateMessageParams{
    ID: req.Id,
    UserID: req.UserId,
    Text: req.Text,
    UpdateUTCNano: updateUTCNano,
  })
  if err != nil {
    logger.Error("message", "failed to update message")
    return nil, err
  }

  return &api.UpdateMessageResponse{
    UpdateUtcNano: updateUTCNano,
  }, nil
}

func (h *Handler) DeleteMessage(ctx context.Context, req *api.DeleteMessageRequest) (
  *api.DeleteMessageResponse, error,
) {
  err := h.controller.DeleteMessage(ctx, logger.Default(), &model.DeleteMessageParams{
    ID: req.Id,
  })
  if err != nil {
    logger.Error("message", "failed to delete message")
    return nil, err
  }

  return &api.DeleteMessageResponse{}, nil
}

func (h *Handler) ReadUserMessages(ctx context.Context, req *api.ReadUserMessagesRequest) (
  *api.ReadUserMessagesResponse, error,
) {
  res, err := h.controller.ReadUserMessages(ctx, logger.Default(), &model.ReadUserMessagesParams{
    UserID:    req.UserId,
    Limit:     req.Limit,
    Offset:    req.Offset,
    Ascending: req.Asc,
  })
  if err != nil {
    logger.Error("message", "failed to read user messages")
    return nil, err
  }

  return &api.ReadUserMessagesResponse{
    Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
    IsLastPage: res.IsLastPage,
  }, nil
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (
  *api.GetServersResponse, error,
) {
  srvs, err := h.controller.GetServers(ctx, logger.Default())
  if err != nil {
    logger.Error("message", "failed to get servers")
    return nil, err
  }

  return &api.GetServersResponse{
    Servers: srvs,
  }, nil
}
