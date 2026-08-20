package orchestrator

import (
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/am"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/messages/internal/domain"
	threadspkg "github.com/bd878/gallery/server/threads/pkg"
)

const CreateMessageReplyChannel = "gallery.messages.replies.CreateMessage"
const CreateMessageSagaName = "messages.CreateMessage"

type orchestratorHandler[T ddd.Event] struct {
	stream am.CommandStream
}

func NewOrchestratorHandlers(stream am.CommandStream) *orchestratorHandler[ddd.Event] {
	return &orchestratorHandler[ddd.Event]{stream}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event]) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		orchestratorEventHandlers := di.Get(ctx, "orchestratorEventHandlers").(ddd.EventHandler[ddd.Event])

		return orchestratorEventHandlers.HandleEvent(ctx, event)
	})

	subscriber.Subscribe(handlers, domain.MessageCreatedEvent)
}

func (h orchestratorHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	switch event.EventName() {
	case domain.MessageCreatedEvent:
		return h.onMessageCreated(ctx, event)
	}
	return nil
}

func (h orchestratorHandler[T]) onMessageCreated(ctx context.Context, event ddd.Event) error {
	slog.Debug("handle event",
		slog.String("name", event.EventName()),
		slog.Any("id", event.ID()),
		slog.Any("payload", event.Payload()),
	)

	payload := event.Payload().(*domain.MessageCreated)
	data, err := proto.Marshal(&threads.CreateThread{
		ThreadId: payload.ID,
		UserId: payload.UserID,
		ParentId: payload.ThreadID,
		Private: &payload.Private,
		Name: payload.Name,
	})
	if err != nil {
		return err
	}

	cmd := am.NewCommand(threadspkg.CreateThreadCommand, threadspkg.CommandChannel, data)

	cmd.Metadata().Set(am.CommandReplyChannelHdr, CreateMessageReplyChannel)
	cmd.Metadata().Set(am.CommandHdrPrefix + "SAGA_ID", cmd.ID())
	cmd.Metadata().Set(am.CommandHdrPrefix + "SAGA_NAME", CreateMessageSagaName)

	return h.stream.Publish(ctx, threadspkg.CommandChannel, cmd)
}
