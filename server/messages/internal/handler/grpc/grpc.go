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
	SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) error
	UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error
	DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error
	ReadMessage(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (*model.Message, error)
	ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (*model.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
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
	logger.Debugw("save message", "text", req.Message)

	req.Message.CreateUtcNano = time.Now().UnixNano()
	req.Message.UpdateUtcNano = time.Now().UnixNano()
	req.Message.Id = utils.RandomID()

	err := h.controller.SaveMessage(ctx, logger.Default(), model.MessageFromProto(req.Message)) 
	if err != nil {
		logger.Errorw("failed to save message", "error", err)
		return nil, err
	}

	logger.Debugw("save message ok", "id", req.Message.Id)
	return &api.SaveMessageResponse{
		Id: req.Message.Id,
		CreateUtcNano: req.Message.CreateUtcNano,
		UpdateUtcNano: req.Message.UpdateUtcNano,
		Private: req.Message.Private,
	}, nil
}

func (h *Handler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (
	*api.UpdateMessageResponse, error,
) {
	logger.Debugw("update message", "id", req.Id)

	updateUTCNano := time.Now().UnixNano()

	err := h.controller.UpdateMessage(ctx, logger.Default(), &model.UpdateMessageParams{
		ID: req.Id,
		UserID: req.UserId,
		Text: req.Text,
		UpdateUTCNano: updateUTCNano,
		Private: req.Private,
	})
	if err != nil {
		logger.Errorw("failed to update message", "error", err)
		return nil, err
	}

	logger.Debugw("update message ok", "id", req.Id)
	return &api.UpdateMessageResponse{
		UpdateUtcNano: updateUTCNano,
	}, nil
}

func (h *Handler) DeleteMessage(ctx context.Context, req *api.DeleteMessageRequest) (
	*api.DeleteMessageResponse, error,
) {
	err := h.controller.DeleteMessage(ctx, logger.Default(), &model.DeleteMessageParams{
		ID: req.Id,
		UserID: req.UserId,
	})
	if err != nil {
		logger.Errorw("failed to delete message", "error", err)
		return nil, err
	}

	return &api.DeleteMessageResponse{}, nil
}

func (h *Handler) ReadThreadMessages(ctx context.Context, req *api.ReadThreadMessagesRequest) (
	*api.ReadThreadMessagesResponse, error,
) {
	logger.Debugw("grpc read thread messages", "user_id", req.UserId, "thread_id", req.ThreadId)
	res, err := h.controller.ReadThreadMessages(ctx, logger.Default(), &model.ReadThreadMessagesParams{
		UserID:    req.UserId,
		ThreadID:  req.ThreadId,
		Limit:     req.Limit,
		Offset:    req.Offset,
		Ascending: req.Asc,
		Private:    req.Private,
	})
	if err != nil {
		logger.Errorw("failed to read thread messages", "error", err)
		return nil, err
	}

	logger.Debugw("grpc read thread messages ok", "thread_id", req.ThreadId, "private", req.Private)
	return &api.ReadThreadMessagesResponse{
		Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, nil
}

func (h *Handler) ReadAllMessages(ctx context.Context, req *api.ReadMessagesRequest) (
	*api.ReadMessagesResponse, error,
) {
	res, err := h.controller.ReadAllMessages(ctx, logger.Default(), &model.ReadMessagesParams{
		UserID:    req.UserId,
		Limit:     req.Limit,
		Offset:    req.Offset,
		Ascending: req.Asc,
		Private:    req.Private,
	})
	if err != nil {
		logger.Errorw("failed to read user messages", "error", err)
		return nil, err
	}

	return &api.ReadMessagesResponse{
		Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, nil
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (
	*api.GetServersResponse, error,
) {
	srvs, err := h.controller.GetServers(ctx, logger.Default())
	if err != nil {
		logger.Errorw("failed to get servers", "error", err)
		return nil, err
	}

	return &api.GetServersResponse{
		Servers: srvs,
	}, nil
}

func (h *Handler) ReadOneMessage(ctx context.Context, req *api.ReadOneMessageRequest) (
	*api.Message, error,
) {
	logger.Debugw("read one message", "user_ids", req.UserIds, "id", req.Id)

	message, err := h.controller.ReadMessage(ctx, logger.Default(), &model.ReadOneMessageParams{
		ID: req.Id,
		UserIDs: req.UserIds,
	})
	if err != nil {
		logger.Errorw("failed to read one message", "user_ids", req.UserIds, "message_id", req.Id)
		return nil, err
	}

	logger.Debugw("read one message ok", "user_ids", req.UserIds, "id", req.Id)
	return model.MessageToProto(message), nil
}