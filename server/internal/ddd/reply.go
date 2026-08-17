package ddd

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type (
	ReplyHandler[T Reply] interface {
		HandleReply(ctx context.Context, reply T) error
	}

	ReplyHandlerFunc[T Reply] func(ctx context.Context, reply T) error

	ReplyOption interface {
		configureReply(*reply)
	}

	ReplyPayload any

	Reply interface {
		IDer
		ReplyName() string
		Payload() ReplyPayload
		Metadata() Metadata
		OccurredAt() time.Time
	}

	reply struct {
		id         string
		name       string
		occurredAt time.Time
		payload    ReplyPayload
		metadata   Metadata
	}
)

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload, options ...ReplyOption) Reply {
	return newReply(name, payload, options...)
}

func newReply(name string, payload ReplyPayload, options ...ReplyOption) reply {
	rep := reply{
		id:         uuid.New().String(),
		name:       name,
		occurredAt: time.Now(),
		payload:    payload,
		metadata:   make(Metadata),
	}

	for _, option := range options {
		option.configureReply(&rep)
	}

	return rep
}

func (e reply) ID() string { return e.id }
func (e reply) ReplyName() string     { return e.name }
func (e reply) Payload() ReplyPayload { return e.payload }
func (e reply) Metadata() Metadata    { return e.metadata }
func (e reply) OccurredAt() time.Time { return e.occurredAt }

func (f ReplyHandlerFunc[T]) HandleReply(ctx context.Context, reply T) error {
	return f(ctx, reply)
}
