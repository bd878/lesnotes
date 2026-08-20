package stream

import (
	"context"

	"log/slog"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/threads"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/threads/pkg"
)

type commandHandlers struct {
	threads ThreadsController
}

func NewCommandHandlers(threads ThreadsController) ddd.CommandHandler[am.CommandMessage] {
	return commandHandlers{threads: threads}
}

func RegisterCommandHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) error {
	return subscriber.Subscribe(pkg.CommandChannel, handlers, am.GroupName("threads-commands"))
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd am.CommandMessage) (ddd.Reply, error) {
	slog.Debug("handle command", slog.String("name", cmd.CommandName()))

	switch cmd.CommandName() {
	case pkg.CreateThreadCommand:
		return h.doCreateThread(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doCreateThread(ctx context.Context, cmd am.CommandMessage) (ddd.Reply, error) {
	m := &threads.CreateThread{}
	if err := proto.Unmarshal(cmd.Data(), m); err != nil {
		return nil, err
	}

	return nil, h.threads.CreateThread(ctx, m.ThreadId, m.UserId, m.ParentId, m.Name, "", "", true)
}