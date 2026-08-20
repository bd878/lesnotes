package ddd

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type (
	CommandHandler[T Command] interface {
		HandleCommand(ctx context.Context, cmd T) (Reply, error)
	}

	CommandOption interface {
		configureCommand(*command)
	}

	Command interface {
		IDer
		CommandName() string
		Data() []byte
		Metadata() Metadata
		OccurredAt() time.Time
	}

	command struct {
		id         string
		name       string
		data []byte
		metadata Metadata
		occurredAt time.Time
	}
)

var _ Command = (*command)(nil)

func NewCommand(name string, data []byte, options ...CommandOption) Command {
	return newCommand(name, data, options...)
}

func newCommand(name string, data []byte, options ...CommandOption) command {
	evt := command{
		id:         uuid.New().String(),
		name:       name,
		data: data,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureCommand(&evt)
	}

	return evt
}

func (e command) ID() string { return e.id }
func (e command) CommandName() string { return e.name }
func (e command) Data() []byte { return e.data }
func (e command) Metadata() Metadata { return e.metadata }
func (e command) OccurredAt() time.Time { return e.occurredAt }
