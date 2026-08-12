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
	EventStream = MessagePublisher[EventMessage]

	EventMessage interface {
		Message
		OccurredAt() time.Time
		Metadata() ddd.Metadata
		Data() []byte
	}

	eventMessage struct {
		metadata ddd.Metadata
		occurredAt time.Time
		id string
		name string
		data []byte
	}

	eventStream struct {
		stream RawMessageStream
	}
)

var _ EventStream = (*eventStream)(nil)
var _ EventMessage = (*eventMessage)(nil)

type eventMsgHandler struct {
	handler MessageHandler[EventMessage]
}

func NewEventStream(stream RawMessageStream) EventStream {
	return eventStream{
		stream: stream,
	}
}

func NewEventMessage(id, name string, data []byte, metadata ddd.Metadata) *eventMessage {
	return &eventMessage{
		id: id,
		name: name,
		data: data,
		metadata: metadata,
		occurredAt: time.Now(),
	}
}

func (m eventMessage) ID() string { return m.id }
func (m eventMessage) MessageName() string { return m.name }
func (m eventMessage) Data() []byte { return m.data }
func (m eventMessage) Metadata() ddd.Metadata { return m.metadata }
func (m eventMessage) OccurredAt() time.Time { return m.occurredAt }

func (s eventStream) Publish(ctx context.Context, topicName string, msg EventMessage) error {
	metadata, err := structpb.NewStruct(msg.Metadata())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&api.EventMessageData{
		Payload:    msg.Data(),
		OccurredAt: timestamppb.New(msg.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, NewRawMessage(msg.ID(), msg.MessageName(), topicName, data))
}

func NewEventMessageHandler(handler MessageHandler[EventMessage]) RawMessageHandler {
	return eventMsgHandler{handler}
}

func (h eventMsgHandler) HandleMessage(ctx context.Context, msg RawMessage) error {
	var eventData api.EventMessageData

	err := proto.Unmarshal(msg.Data(), &eventData)
	if err != nil {
		return err
	}

	eventMsg := eventMessage{
		id: msg.ID(),
		name: msg.MessageName(),
		metadata: eventData.GetMetadata().AsMap(),
		occurredAt: eventData.GetOccurredAt().AsTime(),
		data: eventData.GetPayload(),
	}

	return h.handler.HandleMessage(ctx, eventMsg)
}