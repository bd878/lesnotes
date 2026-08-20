package sec

import (
	"context"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	Orchestrator[T any] interface {
		Start(ctx context.Context, id string, data T) error
		ReplyTopic() string
		HandleReply(ctx context.Context, reply ddd.Reply) error
	}
)