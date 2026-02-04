package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type MessagesController interface {
	SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, userID int64, private bool, name string) (err error)
	UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, userID int64) (err error)
	DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error)
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
	PublishMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error)
	ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *model.Message, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type MessagesHandler struct {
	api.UnimplementedMessagesServer
	controller MessagesController
}

func NewMessagesHandler(ctrl MessagesController) *MessagesHandler {
	handler := &MessagesHandler{controller: ctrl}

	return handler
}

func (h *MessagesHandler) SaveMessage(ctx context.Context, req *api.SaveMessageRequest) (resp *api.SaveMessageResponse, err error) {
	err = h.controller.SaveMessage(ctx, req.Id, req.Text, req.Title, req.FileIds, req.UserId, req.Private, req.Name)

	resp = &api.SaveMessageResponse{}

	return
}

func (h *MessagesHandler) DeleteUserMessages(ctx context.Context, req *api.DeleteUserMessagesRequest) (resp *api.DeleteUserMessagesResponse, err error) {
	err = h.controller.DeleteUserMessages(ctx, req.UserId)

	resp = &api.DeleteUserMessagesResponse{}

	return
}

func (h *MessagesHandler) UpdateMessage(ctx context.Context, req *api.UpdateMessageRequest) (resp *api.UpdateMessageResponse, err error) {
	err = h.controller.UpdateMessage(ctx, req.Id, req.Text, req.Title, req.Name, req.FileIds, req.UserId)

	resp = &api.UpdateMessageResponse{}

	return
}

func (h *MessagesHandler) DeleteMessages(ctx context.Context, req *api.DeleteMessagesRequest) (resp *api.DeleteMessagesResponse, err error) {
	err = h.controller.DeleteMessages(ctx, req.Ids, req.UserId)

	resp = &api.DeleteMessagesResponse{}

	return
}

func (h *MessagesHandler) PublishMessages(ctx context.Context, req *api.PublishMessagesRequest) (resp *api.PublishMessagesResponse, err error) {
	err = h.controller.PublishMessages(ctx, req.Ids, req.UserId)

	resp = &api.PublishMessagesResponse{}

	return
}

func (h *MessagesHandler) PrivateMessages(ctx context.Context, req *api.PrivateMessagesRequest) (resp *api.PrivateMessagesResponse, err error) {
	err = h.controller.PrivateMessages(ctx, req.Ids, req.UserId)

	resp = &api.PrivateMessagesResponse{}

	return
}

func (h *MessagesHandler) ReadBatchMessages(ctx context.Context, req *api.ReadBatchMessagesRequest) (resp *api.ReadBatchMessagesResponse, err error) {
	list, err := h.controller.ReadBatchMessages(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadBatchMessagesResponse{
		Messages:   model.MapMessagesToProto(model.MessageToProto, list),
	}

	return
}

func (h *MessagesHandler) ReadMessages(ctx context.Context, req *api.ReadMessagesRequest) (resp *api.ReadMessagesResponse, err error) {
	list, isLastPage, err := h.controller.ReadMessages(ctx, req.UserId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadMessagesResponse{
		Messages:   model.MapMessagesToProto(model.MessageToProto, list),
		IsLastPage: isLastPage,
	}

	return
}

func (h *MessagesHandler) GetServers(ctx context.Context, _ *api.GetServersRequest) (resp *api.GetServersResponse, err error) {
	servers, err := h.controller.GetServers(ctx)
	if err != nil {
		return nil, err
	}

	resp = &api.GetServersResponse{
		Servers: servers,
	}

	return
}

func (h *MessagesHandler) ReadMessage(ctx context.Context, req *api.ReadMessageRequest) (resp *api.Message, err error) {
	message, err := h.controller.ReadMessage(ctx, req.Id, req.Name, req.UserIds)
	if err != nil {
		return nil, err
	}

	resp = model.MessageToProto(message)

	return
}
