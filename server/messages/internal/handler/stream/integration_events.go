package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/api"
	users "github.com/bd878/gallery/server/users/pkg/events"
)

type MessagesController interface {
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
}

type integrationHandlers struct {
	messages MessagesController
}

var _ am.RawMessageHandler = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController) am.RawMessageHandler {
	return integrationHandlers{
		messages: messages,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) (err error) {
	err = subscriber.Subscribe(users.UsersChannel, handlers, am.GroupName("messages-users"))
	if err != nil {
		return
	}

	return
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.IncomingMessage) error {
	logger.Debugw("handle message", "name", msg.MessageName(), "subject", msg.Subject())

	switch msg.MessageName() {
	case users.UserDeletedEvent:
		return h.handleUserDeleted(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleUserDeleted(ctx context.Context, msg am.IncomingMessage) error {
	m := &api.UserDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.DeleteUserMessages(ctx, m.GetUserId())
}
