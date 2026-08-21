package tm

import (
	"fmt"
	"log/slog"
	"errors"
	"context"
	"github.com/bd878/gallery/server/internal/am"
)

type ErrDuplicateMessage string

type InboxStore interface {
	Save(ctx context.Context, msg am.RawMessage) error
}

type inbox struct {
	handler am.RawMessageHandler
	store InboxStore
}

var _ am.RawMessageHandler = (*inbox)(nil)

func NewInboxHandlerMiddleware(store InboxStore) am.RawMessageHandlerMiddleware {
	i := inbox{store: store}

	return func(handler am.RawMessageHandler) am.RawMessageHandler {
		i.handler = handler

		return i
	}
}

func (i inbox) HandleMessage(ctx context.Context, msg am.RawMessage) error {
	err := i.store.Save(ctx, msg)
	if err != nil {
		var errDupe ErrDuplicateMessage
		if errors.Is(err, errDupe) {
			slog.Debug("duplicate message",
				slog.String("subject", msg.Subject()),
				slog.String("name", msg.MessageName()),
			)
			return nil
		}
		return err
	}

	return i.handler.HandleMessage(ctx, msg)
}

func (e ErrDuplicateMessage) Error() string {
	return fmt.Sprintf("duplicate message id encountered: %s", string(e))
}