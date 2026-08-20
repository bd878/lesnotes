package am

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	ReplyMessage interface {
		ddd.Reply
	}

	ReplyPublisher = MessagePublisher[ReplyMessage]
	ReplySubscriber = MessageSubscriber[ReplyMessage]
	ReplyStream     = MessageStream[ReplyMessage, ReplyMessage]

	replyStream struct {
		stream RawMessageStream
	}

	replyMessage struct {
		id string
		name string
		data []byte
		metadata ddd.Metadata
		occurredAt time.Time
	}
)

var _ ReplyMessage = (*replyMessage)(nil)

var _ ReplyStream = (*replyStream)(nil)

func NewReplyStream(stream RawMessageStream) ReplyStream {
	return replyStream{
		stream: stream,
	}
}

func (s replyStream) Publish(ctx context.Context, topicName string, reply ReplyMessage) error {
	metadata, err := structpb.NewStruct(reply.Metadata())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&api.ReplyMessageData{
		Payload:    reply.Data(),
		OccurredAt: timestamppb.New(reply.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id: reply.ID(),
		name: reply.ReplyName(),
		subject: topicName,
		data: data,
	})
}

func (s replyStream) Subscribe(topicName string, handler MessageHandler[ReplyMessage], options ...SubscriberOption) error {
	fn := RawMessageHandlerFunc(func(ctx context.Context, msg RawMessage) error {
		var replyData api.ReplyMessageData

		err := proto.Unmarshal(msg.Data(), &replyData)
		if err != nil {
			return err
		}

		replyMsg := replyMessage{
			id: msg.ID(),
			name: msg.MessageName(),
			data: replyData.GetPayload(),
			metadata: replyData.GetMetadata().AsMap(),
			occurredAt: replyData.GetOccurredAt().AsTime(),
		}

		return handler.HandleMessage(ctx, replyMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (m replyMessage) ID() string { return m.id }
func (m replyMessage) ReplyName() string { return m.name }
func (m replyMessage) Data() []byte { return m.data }
func (m replyMessage) Metadata() ddd.Metadata { return m.metadata }
func (m replyMessage) OccurredAt() time.Time { return m.occurredAt }

type replyMsgHandler struct {
	handler ddd.ReplyHandler[ddd.Reply]
}

func NewReplyMessageHandler(handler ddd.ReplyHandler[ddd.Reply]) RawMessageHandler {
	return replyMsgHandler{
		handler: handler,
	}
}

func (h replyMsgHandler) HandleMessage(ctx context.Context, msg RawMessage) error {
	var replyData api.ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {
		return err
	}

	replyMsg := replyMessage{
		id: msg.ID(),
		name: msg.MessageName(),
		data: replyData.GetPayload(),
		occurredAt: replyData.GetOccurredAt().AsTime(),
		metadata: replyData.GetMetadata().AsMap(),
	}

	return h.handler.HandleReply(ctx, replyMsg)
}