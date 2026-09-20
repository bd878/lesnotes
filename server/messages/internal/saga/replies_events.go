package saga

import (
	"context"
	"log/slog"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/sec"
)

const messageReplyChannel = "gallery.messages.replies"

type replyHandlers struct {
	createMessageSaga sec.Orchestrator
	deleteMessageSaga sec.Orchestrator
}

var _ ddd.ReplyHandler[am.ReplyMessage] = (*replyHandlers)(nil)

func NewReplyEventHandlers(createMessageSaga, deleteMessageSaga sec.Orchestrator) ddd.ReplyHandler[am.ReplyMessage] {
	return replyHandlers{createMessageSaga, deleteMessageSaga}
}

func (h replyHandlers) HandleReply(ctx context.Context, reply am.ReplyMessage) error {
	slog.Debug("handle reply",
		slog.String("reply_name", reply.ReplyName()),
		slog.String("subject", reply.Subject()),
		slog.String("id", reply.ID()),
	)

	switch reply.Subject() {
	case h.createMessageSaga.ReplyTopic():
		return h.createMessageSaga.HandleReply(ctx, reply)
	case h.deleteMessageSaga.ReplyTopic():
		return h.deleteMessageSaga.HandleReply(ctx, reply)
	default:
		slog.Error("unknown reply name",
			slog.String("reply_name", reply.ReplyName()),
			slog.String("subject", reply.Subject()),
		)
	}

	return nil
}