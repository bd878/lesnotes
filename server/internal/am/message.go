package am

import (
	"time"
	"context"

	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	MessageBase interface {
		ID() string
		Subject() string
		MessageName() string
		Metadata() ddd.Metadata
		SentAt() time.Time
	}

	Message interface {
		MessageBase
		Data() []byte
	}

	IncomingMessage interface {
		MessageBase
		Data() []byte // TODO: remove, use ddd.Event in integrationEvents
	}

	MessageHandler interface {
		HandleMessage(ctx context.Context, msg IncomingMessage) error
	}

	MessageHandlerFunc func(ctx context.Context, msg IncomingMessage) error
	MessagePublisherFunc func(ctx context.Context, topicName string, msg Message) error

	MessagePublisher interface {
		Publish(ctx context.Context, topicName string, msg Message) error
	}

	MessagePublisherMiddleware = func(next MessagePublisher) MessagePublisher
	MessageHandlerMiddleware = func(next MessageHandler) MessageHandler

	MessageSubscriber interface {
		Subscribe(topicName string, handler MessageHandler) error
	}

	MessageStream interface {
		MessagePublisher
		MessageSubscriber
	}

	messagePublisher struct {
		publisher MessagePublisher
	}

	messageSubscriber struct {
		subscriber MessageSubscriber
		mws        []MessageHandlerMiddleware
	}
)

func (f MessagePublisherFunc) Publish(ctx context.Context, topicName string, msg Message) error {
	return f(ctx, topicName, msg)
}

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	return f(ctx, msg)
}

func NewMessagePublisher(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return messagePublisher{
		publisher: MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (p messagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	return p.publisher.Publish(ctx, topicName, msg)
}

func NewMessageSubscriber(subscriber MessageSubscriber, mws ...MessageHandlerMiddleware) MessageSubscriber {
	return messageSubscriber{
		subscriber: subscriber,
		mws:        mws,
	}
}

func (s messageSubscriber) Subscribe(topicName string, handler MessageHandler) error {
	return s.subscriber.Subscribe(topicName, MessageHandlerWithMiddleware(handler, s.mws...))
}
