package am

import (
	"context"
)

type (
	Message interface {
		ID()          string
		Subject()     string
		MessageName() string
	}

	IncomingMessage interface {
		Message
		Ack()   error
		Nak()  error
	}

	MessageHandler[I IncomingMessage] interface {
		HandleMessage(ctx context.Context, msg I) error
	}

	MessageHandlerFunc[I IncomingMessage] func(ctx context.Context, msg I) error

	MessagePublisher[O any] interface {
		Publish(ctx context.Context, topicName string, v O) error
	}

	MessageSubscriber[I IncomingMessage] interface {
		Subscribe(topicName string, handler MessageHandler[I]) error
	}

	MessageStream[O any, I IncomingMessage] interface {
		MessagePublisher[O]
		MessageSubscriber[I]
	}
)

func (f MessageHandlerFunc[I]) HandleMessage(ctx context.Context, msg I) error {
	return f(ctx, msg)
}