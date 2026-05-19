package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Stream struct {
	nc  *nats.Conn
}

var _ am.MessageStream = (*Stream)(nil)

func NewStream(nc *nats.Conn) *Stream {
	return &Stream{nc: nc}
}

func (s *Stream) Publish(ctx context.Context, topicName string, message am.Message) error {
	metadata, err := structpb.NewStruct(message.Metadata())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&api.StreamMessage{
		Id:   message.ID(),
		Name: message.MessageName(),
		Data: message.Data(),
		Metadata: metadata,
		SentAt: timestamppb.New(message.SentAt()),
	})
	if err != nil {
		return err
	}
	return s.nc.Publish(topicName, data)
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler) error {
	_, err := s.nc.Subscribe(topicName, s.handleMsg(topicName, handler))

	return err
}

func (s *Stream) handleMsg(topicName string, handler am.MessageHandler) func(*nats.Msg) {
	return func(natsMsg *nats.Msg) {
		var err error

		m := &api.StreamMessage{}
		err = proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			logger.Errorw("failed ot unmarshal nats message", "error", err)
			return
		}

		msg := &rawMessage{
			id:      m.GetId(),
			name:    m.GetName(),
			data:    m.GetData(),
			subject: topicName,
			metadata: m.GetMetadata().AsMap(),
			sentAt:  m.SentAt.AsTime(),
		}

		wCtx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		select {
		case err = <-errc:
			if err != nil {
				logger.Errorw("error while handling message", "error", err)
			}
			return
		case <-wCtx.Done():
			return
		}
	}
}