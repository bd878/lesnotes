package grpc

import (
	"time"
	"context"

	"github.com/google/uuid"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type Controller interface {
	SaveMessage(ctx context.Context, message *model.Message) error
	UpdateMessage(ctx context.Context, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, params *model.DeleteMessageParams) error
	DeleteAllUserMessages(ctx context.Context, params *model.DeleteAllUserMessagesParams) error
	DeleteMessages(ctx context.Context, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error)
	PublishMessages(ctx context.Context, params *model.PublishMessagesParams) (*model.PublishMessagesResult, error)
	PrivateMessages(ctx context.Context, params *model.PrivateMessagesParams) (*model.PrivateMessagesResult, error)
	ReadMessage(ctx context.Context, id int64, userIDs []int64) (*model.Message, error)
	ReadAllMessages(ctx context.Context, params *model.ReadMessagesParams) (messages []*model.Message, isLastPage bool, err error,)
	ReadThreadMessages(ctx context.Context, params *model.ReadThreadMessagesParams) (messages []*model.Message, isLastPage bool, err error,)
	GetServers(ctx context.Context) ([]*api.Server, error)
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
	req.Message.Id = int64(utils.RandomID())
	req.Message.Name = uuid.New().String()

	err := h.controller.SaveMessage(ctx, model.MessageFromProto(req.Message)) 
	if err != nil {
		return nil, err
	}

	return &api.SaveMessageResponse{
		Id: req.Message.Id,
		CreateUtcNano: req.Message.CreateUtcNano,
		UpdateUtcNano: req.Message.UpdateUtcNano,
		Private: req.Message.Private,
		Name: req.Message.Name,
	}, nil
}

func (h *Handler) DeleteAllUserMessages(ctx context.Context, req *api.DeleteAllUserMessagesRequest) (*api.DeleteAllUserMessagesResponse, error) {
	err := h.controller.DeleteAllUserMessages(ctx, &model.DeleteAllUserMessagesParams{
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &api.DeleteAllUserMessagesResponse{}, nil
}

func (h *Handler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (*api.UpdateMessageResponse, error) {
	updateUTCNano := time.Now().UnixNano()

	res, err := h.controller.UpdateMessage(ctx, &model.UpdateMessageParams{
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
	err := h.controller.DeleteMessage(ctx, &model.DeleteMessageParams{
		ID: req.Id,
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &api.DeleteMessageResponse{}, nil
}

func (h *Handler) DeleteMessages(ctx context.Context, req *api.DeleteMessagesRequest) (*api.DeleteMessagesResponse, error) {
	res, err := h.controller.DeleteMessages(ctx, &model.DeleteMessagesParams{
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
	_, err := h.controller.PublishMessages(ctx, &model.PublishMessagesParams{
		IDs: req.Ids,
		UserID: req.UserId,
		UpdateUTCNano: updateUTCNano,
	})

	return &api.PublishMessagesResponse{UpdateUtcNano: updateUTCNano}, err
}

func (h *Handler) PrivateMessages(ctx context.Context, req *api.PrivateMessagesRequest) (*api.PrivateMessagesResponse, error) {
	updateUTCNano := time.Now().UnixNano()
	_, err := h.controller.PrivateMessages(ctx, &model.PrivateMessagesParams{
		IDs: req.Ids,
		UserID: req.UserId,
		UpdateUTCNano: updateUTCNano,
	})

	return &api.PrivateMessagesResponse{UpdateUtcNano: updateUTCNano}, err
}

func (h *Handler) ReadThreadMessages(ctx context.Context, req *api.ReadThreadMessagesRequest) (*api.ReadThreadMessagesResponse, error) {
	messages, isLastPage, err := h.controller.ReadThreadMessages(ctx, &model.ReadThreadMessagesParams{
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
		Messages: model.MapMessagesToProto(model.MessageToProto, messages),
		IsLastPage: isLastPage,
	}, nil
}

func (h *Handler) ReadAllMessages(ctx context.Context, req *api.ReadMessagesRequest) (*api.ReadMessagesResponse, error) {
	messages, isLastPage, err := h.controller.ReadAllMessages(ctx, &model.ReadMessagesParams{
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
		Messages: model.MapMessagesToProto(model.MessageToProto, messages),
		IsLastPage: isLastPage,
	}, nil
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (*api.GetServersResponse, error) {
	srvs, err := h.controller.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetServersResponse{
		Servers: srvs,
	}, nil
}

func (h *Handler) ReadOneMessage(ctx context.Context, req *api.ReadMessageRequest) (*api.Message, error) {
	message, err := h.controller.ReadMessage(ctx, req.Id, req.UserIds)
	if err != nil {
		return nil, err
	}

	return model.MessageToProto(message), nil
}