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

	ReplyOption interface {
		configureReply(*reply)
	}

	Reply interface {
		IDer
		ReplyName() string
		Data() []byte
		Metadata() Metadata
		OccurredAt() time.Time
	}

	reply struct {
		id         string
		name       string
		occurredAt time.Time
		data []byte
		metadata   Metadata
	}
)

var _ Reply = (*reply)(nil)

func NewReply(name string, data []byte, options ...ReplyOption) Reply {
	return newReply(name, data, options...)
}

func newReply(name string, data []byte, options ...ReplyOption) reply {
	rep := reply{
		id:         uuid.New().String(),
		name:       name,
		occurredAt: time.Now(),
		data: data,
		metadata:   make(Metadata),
	}

	for _, option := range options {
		option.configureReply(&rep)
	}

	return rep
}

func (e reply) ID() string { return e.id }
func (e reply) ReplyName() string     { return e.name }
func (e reply) Data() []byte { return e.data }
func (e reply) Metadata() Metadata    { return e.metadata }
func (e reply) OccurredAt() time.Time { return e.occurredAt }
