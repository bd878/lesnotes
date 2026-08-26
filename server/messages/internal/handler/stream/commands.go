package stream

import (
	"context"
	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/messages/pkg"
)

type commandHandlers struct {
	messages MessagesController
}

func NewCommandHandlers(messages MessagesController) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{messages: messages}
}

func RegisterCommandHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) error {
	return subscriber.Subscribe(pkg.CommandChannel, handlers, am.GroupName("messages-commands"))
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	slog.Debug("handle command", slog.String("name", cmd.CommandName()))

	switch cmd.CommandName() {
	case pkg.DeleteMessageCommand:
		return h.doDeleteMessage(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doDeleteMessage(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	m := &messages.DeleteMessage{}
	if err := proto.Unmarshal(cmd.Data(), m); err != nil {
		return nil, err
	}

	return nil, h.messages.DeleteMessages(ctx, []int64{m.Id}, m.UserId)
}
