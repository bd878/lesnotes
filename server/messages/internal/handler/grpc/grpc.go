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
	UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error
	DeleteMessages(ctx context.Context, log *logger.Logger, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error)
	PublishMessages(ctx context.Context, log *logger.Logger, params *model.PublishMessagesParams) (*model.PublishMessagesResult, error)
	PrivateMessages(ctx context.Context, log *logger.Logger, params *model.PrivateMessagesParams) (*model.PrivateMessagesResult, error)
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

func (h *Handler) SaveMessage(ctx context.Context, req *api.SaveMessageRequest) (*api.SaveMessageResponse, error) {
	req.Message.CreateUtcNano = time.Now().UnixNano()
	req.Message.UpdateUtcNano = time.Now().UnixNano()
	req.Message.Id = utils.RandomID()

	err := h.controller.SaveMessage(ctx, logger.Default(), model.MessageFromProto(req.Message)) 
	if err != nil {
		return nil, err
	}

	return &api.SaveMessageResponse{
		Id: req.Message.Id,
		CreateUtcNano: req.Message.CreateUtcNano,
		UpdateUtcNano: req.Message.UpdateUtcNano,
		Private: req.Message.Private,
	}, nil
}

func (h *Handler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (*api.UpdateMessageResponse, error) {
	updateUTCNano := time.Now().UnixNano()

	res, err := h.controller.UpdateMessage(ctx, logger.Default(), &model.UpdateMessageParams{
		ID: req.Id,
		UserID: req.UserId,
		ThreadID: req.ThreadId,
		Text: req.Text,
		UpdateUTCNano: updateUTCNano,
		Private: req.Private,
	})
	if err != nil {
		return nil, err
	}

	return &api.UpdateMessageResponse{
		UpdateUtcNano: updateUTCNano,
		Private: res.Private,
	}, nil
}

func (h *Handler) DeleteMessage(ctx context.Context, req *api.DeleteMessageRequest) (*api.DeleteMessageResponse, error) {
	err := h.controller.DeleteMessage(ctx, logger.Default(), &model.DeleteMessageParams{
		ID: req.Id,
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &api.DeleteMessageResponse{}, nil
}

func (h *Handler) DeleteMessages(ctx context.Context, req *api.DeleteMessagesRequest) (*api.DeleteMessagesResponse, error) {
	res, err := h.controller.DeleteMessages(ctx, logger.Default(), &model.DeleteMessagesParams{
		IDs: req.Ids,
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	ids := make([]*api.DeleteMessageStatus, 0, len(res.IDs))
	for _, status := range res.IDs {
		ids = append(ids, &api.DeleteMessageStatus{Id: status.ID, Ok: status.OK, Explain: status.Explain})
	}
	return &api.DeleteMessagesResponse{Ids: ids}, nil
}

func (h *Handler) PublishMessages(ctx context.Context, req *api.PublishMessagesRequest) (*api.PublishMessagesResponse, error) {
	updateUTCNano := time.Now().UnixNano()
	_, err := h.controller.PublishMessages(ctx, logger.Default(), &model.PublishMessagesParams{
		IDs: req.Ids,
		UserID: req.UserId,
		UpdateUTCNano: updateUTCNano,
	})

	return &api.PublishMessagesResponse{UpdateUtcNano: updateUTCNano}, err
}

func (h *Handler) PrivateMessages(ctx context.Context, req *api.PrivateMessagesRequest) (*api.PrivateMessagesResponse, error) {
	updateUTCNano := time.Now().UnixNano()
	_, err := h.controller.PrivateMessages(ctx, logger.Default(), &model.PrivateMessagesParams{
		IDs: req.Ids,
		UserID: req.UserId,
		UpdateUTCNano: updateUTCNano,
	})

	return &api.PrivateMessagesResponse{UpdateUtcNano: updateUTCNano}, err
}

func (h *Handler) ReadThreadMessages(ctx context.Context, req *api.ReadThreadMessagesRequest) (*api.ReadThreadMessagesResponse, error) {
	res, err := h.controller.ReadThreadMessages(ctx, logger.Default(), &model.ReadThreadMessagesParams{
		UserID:    req.UserId,
		ThreadID:  req.ThreadId,
		Limit:     req.Limit,
		Offset:    req.Offset,
		Ascending: req.Asc,
		Private:    req.Private,
	})
	if err != nil {
		return nil, err
	}

	return &api.ReadThreadMessagesResponse{
		Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, nil
}

func (h *Handler) ReadAllMessages(ctx context.Context, req *api.ReadMessagesRequest) (*api.ReadMessagesResponse, error) {
	res, err := h.controller.ReadAllMessages(ctx, logger.Default(), &model.ReadMessagesParams{
		UserID:    req.UserId,
		Limit:     req.Limit,
		Offset:    req.Offset,
		Ascending: req.Asc,
		Private:    req.Private,
	})
	if err != nil {
		return nil, err
	}

	return &api.ReadMessagesResponse{
		Messages: model.MapMessagesToProto(model.MessageToProto, res.Messages),
		IsLastPage: res.IsLastPage,
	}, nil
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (*api.GetServersResponse, error) {
	srvs, err := h.controller.GetServers(ctx, logger.Default())
	if err != nil {
		return nil, err
	}

	return &api.GetServersResponse{
		Servers: srvs,
	}, nil
}

func (h *Handler) ReadOneMessage(ctx context.Context, req *api.ReadOneMessageRequest) (*api.Message, error) {
	message, err := h.controller.ReadMessage(ctx, logger.Default(), &model.ReadOneMessageParams{
		ID: req.Id,
		UserIDs: req.UserIds,
	})
	if err != nil {
		return nil, err
	}

	return model.MessageToProto(message), nil
}