package am

import (
	"time"
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	CommandMessage interface {
		ddd.Command
	}

	CommandStream interface {
		MessagePublisher[ddd.Command]
		MessageSubscriber[CommandMessage]
	}

	commandStream struct {
		stream RawMessageStream
	}

	commandMessage struct {
		id string
		name string
		data []byte
		metadata ddd.Metadata
		occurredAt time.Time
	}
)

var _ CommandMessage = (*commandMessage)(nil)

var _ CommandStream = (*commandStream)(nil)

func NewCommandStream(stream RawMessageStream) CommandStream {
	return commandStream{
		stream: stream,
	}
}

func (s commandStream) Publish(ctx context.Context, topicName string, msg ddd.Command) error {
	metadata, err := structpb.NewStruct(msg.Metadata())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&api.CommandMessageData{
		Payload:    msg.Data(),
		OccurredAt: timestamppb.New(msg.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id: msg.ID(),
		name: msg.CommandName(),
		subject: topicName,
		data: data,
	})
}

func (s commandStream) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
	fn := MessageHandlerFunc[RawMessage](func(ctx context.Context, msg RawMessage) error {
		var commandData api.CommandMessageData

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		commandMsg := commandMessage{
			id:         msg.ID(),
			name:       msg.MessageName(),
			data:       commandData.GetPayload(),
			metadata:   commandData.GetMetadata().AsMap(),
			occurredAt: commandData.GetOccurredAt().AsTime(),
		}

		return handler.HandleMessage(ctx, commandMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Data() []byte { return c.data }
func (c commandMessage) Metadata() ddd.Metadata      { return c.metadata }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
