package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"log/slog"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/users/pkg/events"
)

type MessagesController interface {
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
}

type integrationHandlers struct {
	messages MessagesController
}

var _ am.MessageHandler[am.EventMessage] = (*integrationHandlers)(nil)

func NewIntegrationEventHandlers(messages MessagesController) am.MessageHandler[am.EventMessage] {
	return integrationHandlers{
		messages: messages,
	}
}

func RegisterIntegrationEventHandlers(stream am.RawMessageStream, handlers am.MessageHandler[am.EventMessage]) (err error) {
	evtMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.RawMessage) error {
		// TODO: open/commit/rollback tx
		evtHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(
				handlers,
			),
			// TODO: inboxMiddleware.(am.RawMessageHandlerMiddleware)
		)

		return evtHandlers.HandleMessage(ctx, msg)
	})

	return stream.Subscribe(events.UsersChannel, evtMsgHandler, am.GroupName("messages-users"))
}

func (h integrationHandlers) HandleMessage(ctx context.Context, msg am.EventMessage) error {
	slog.Debug("handle message", slog.String("name", msg.MessageName()))

	switch msg.MessageName() {
	case events.UserDeletedEvent:
		return h.handleUserDeleted(ctx, msg)
	}

	return nil
}

func (h integrationHandlers) handleUserDeleted(ctx context.Context, msg am.EventMessage) error {
	m := &users.UserDeleted{}
	if err := proto.Unmarshal(msg.Data(), m); err != nil {
		return err
	}

	return h.messages.DeleteUserMessages(ctx, m.GetUserId())
}
