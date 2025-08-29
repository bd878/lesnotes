package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
)

type Controller interface {
	SaveMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (err error)
	UpdateMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private int32) (err error)
	DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error)
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
	PublishMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error)
	ReadMessage(ctx context.Context, id int64, userIDs []int64) (message *messages.Message, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*messages.Message, isLastPage bool, err error)
	ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (messages []*messages.Message, isLastPage bool, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type Handler struct {
	api.UnimplementedMessagesServer
	controller Controller
}

func New(ctrl Controller) *Handler {
	handler := &Handler{controller: ctrl}

	return handler
}

func (h *Handler) SaveMessage(ctx context.Context, req *api.SaveMessageRequest) (resp *api.SaveMessageResponse, err error) {
	err = h.controller.SaveMessage(ctx, req.Id, req.Text, req.FileIds, req.ThreadId, req.UserId, req.Private, req.Name)

	resp = &api.SaveMessageResponse{}

	return
}

func (h *Handler) DeleteUserMessages(ctx context.Context, req *api.DeleteUserMessagesRequest) (resp *api.DeleteUserMessagesResponse, err error) {
	err = h.controller.DeleteUserMessages(ctx, req.UserId)

	resp = &api.DeleteUserMessagesResponse{}

	return
}

func (h *Handler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (resp *api.UpdateMessageResponse, err error) {
	err = h.controller.UpdateMessage(ctx, req.Id, req.Text, nil, req.ThreadId, req.UserId, req.Private)

	resp = &api.UpdateMessageResponse{}

	return
}

func (h *Handler) DeleteMessages(ctx context.Context, req *api.DeleteMessagesRequest) (resp *api.DeleteMessagesResponse, err error) {
	err = h.controller.DeleteMessages(ctx, req.Ids, req.UserId)

	resp = &api.DeleteMessagesResponse{}

	return
}

func (h *Handler) PublishMessages(ctx context.Context, req *api.PublishMessagesRequest) (resp *api.PublishMessagesResponse, err error) {
	err = h.controller.PublishMessages(ctx, req.Ids, req.UserId)

	resp = &api.PublishMessagesResponse{}

	return
}

func (h *Handler) PrivateMessages(ctx context.Context, req *api.PrivateMessagesRequest) (resp *api.PrivateMessagesResponse, err error) {
	err = h.controller.PrivateMessages(ctx, req.Ids, req.UserId)

	resp = &api.PrivateMessagesResponse{}

	return
}

func (h *Handler) ReadThreadMessages(ctx context.Context, req *api.ReadThreadMessagesRequest) (resp *api.ReadThreadMessagesResponse, err error) {
	list, isLastPage, err := h.controller.ReadThreadMessages(ctx, req.UserId, req.ThreadId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadThreadMessagesResponse{
		Messages:   messages.MapMessagesToProto(messages.MessageToProto, list),
		IsLastPage: isLastPage,
	}

	return
}

func (h *Handler) ReadMessages(ctx context.Context, req *api.ReadMessagesRequest) (resp *api.ReadMessagesResponse, err error) {
	list, isLastPage, err := h.controller.ReadMessages(ctx, req.UserId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadMessagesResponse{
		Messages:   messages.MapMessagesToProto(messages.MessageToProto, list),
		IsLastPage: isLastPage,
	}

	return resp, nil
}

func (h *Handler) GetServers(ctx context.Context, _ *api.GetServersRequest) (resp *api.GetServersResponse, err error) {
	servers, err := h.controller.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	resp = &api.GetServersResponse{
		Servers: servers,
	}

	return
}

func (h *Handler) ReadMessage(ctx context.Context, req *api.ReadMessageRequest) (resp *api.Message, err error) {
	message, err := h.controller.ReadMessage(ctx, req.Id, req.UserIds)
	if err != nil {
		return nil, err
	}

	resp = messages.MessageToProto(message)

	return
}