package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/internal/controller"
	"github.com/bd878/gallery/server/messages/internal/machine"
)

type MessagesController interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *api.Message, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*api.Message, isLastPage bool, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*api.Message, err error)
	GetServers(ctx context.Context) (servers []*api.Server, err error)
}

type MessagesHandler struct {
	api.UnimplementedMessagesServer
	api.UnimplementedDistributedServer
	controller MessagesController
}

func NewMessagesHandler(ctrl MessagesController) *MessagesHandler {
	handler := &MessagesHandler{controller: ctrl}

	return handler
}

func (h *MessagesHandler) Apply(ctx context.Context, req *api.Command) (resp *api.CommandResponse, err error) {
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	err = h.controller.Apply(ctx, machine.RequestType(req.ReqType), req.Cmd, duration)
	return
}

func (h *MessagesHandler) ReadBatchMessages(ctx context.Context, req *api.ReadBatchMessagesRequest) (resp *api.ReadBatchMessagesResponse, err error) {
	list, err := h.controller.ReadBatchMessages(ctx, req.UserId, req.Ids)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadBatchMessagesResponse{
		Messages: list,
	}

	return
}

func (h *MessagesHandler) ReadMessages(ctx context.Context, req *api.ReadMessagesRequest) (resp *api.ReadMessagesResponse, err error) {
	list, isLastPage, err := h.controller.ReadMessages(ctx, req.UserId, req.Limit, req.Offset, req.Asc)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadMessagesResponse{
		Messages:   list,
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

func (h *MessagesHandler) ReadMessage(ctx context.Context, req *api.ReadMessageRequest) (resp *api.ReadMessageResponse, err error) {
	message, err := h.controller.ReadMessage(ctx, req.Id, req.Name, req.UserIds)
	if err != nil {
		if errors.Is(err, controller.ErrMessageIsRoot) {
			return &api.ReadMessageResponse{IsRoot: true}, nil
		}

		return nil, err
	}

	return &api.ReadMessageResponse{Message: message}, nil
}
