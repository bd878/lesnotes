package am

import (
	"context"
)

type (
	Message interface {
		ID() string
		MessageName() string
	}

	MessageHandler[I any] interface {
		HandleMessage(ctx context.Context, msg I) error
	}

	MessagePublisher[O any] interface {
		Publish(ctx context.Context, topicName string, msg O) error
	}

	MessageSubscriber[I any] interface {
		Subscribe(topicName string, handler MessageHandler[I], options ...SubscriberOption) error
	}

	MessageStream[O any, I any] interface {
		MessagePublisher[O]
		MessageSubscriber[I]
	}
)
