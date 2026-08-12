package am

import (
	"context"
)

type (
	RawMessageHandler = MessageHandler[RawMessage]

	RawMessageSubscriber = MessageSubscriber[RawMessage]

	RawMessageHandlerFunc       func(ctx context.Context, msg RawMessage) error
	RawMessageStream = MessageStream[RawMessage, RawMessage]
	RawMessageStreamMiddleware = func(stream RawMessageStream) RawMessageStream
	RawMessageHandlerMiddleware = func(handler RawMessageHandler) RawMessageHandler

	RawMessage interface {
		Message
		Data() []byte
		Subject() string
	}

	rawMessage struct {
		id    string
		name  string
		subject string
		data  []byte
	}
)

var _ RawMessage = (*rawMessage)(nil)

func NewRawMessage(id, name, subject string, data []byte) *rawMessage {
	return &rawMessage{
		id: id,
		name: name,
		subject: subject,
		data: data,
	}
}

func (m rawMessage) ID() string { return m.id }
func (m rawMessage) Subject() string { return m.subject }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte { return m.data }

func (f RawMessageHandlerFunc) HandleMessage(ctx context.Context, cmd RawMessage) error {
	return f(ctx, cmd)
}

func RawMessageStreamWithMiddleware(stream RawMessageStream, mws ...RawMessageStreamMiddleware) RawMessageStream {
	s := stream
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}
	return s
}

func RawMessageHandlerWithMiddleware(handler RawMessageHandler, mws ...RawMessageHandlerMiddleware) RawMessageHandler {
	h := handler
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
