package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Stream struct {
	nc  *nats.Conn
}

var _ am.RawMessageStream = (*Stream)(nil)

func NewStream(nc *nats.Conn) *Stream {
	return &Stream{nc: nc}
}

func (s *Stream) Publish(ctx context.Context, topicName string, message am.Message) error {
	data, err := proto.Marshal(&api.StreamMessage{
		Id:   message.ID(),
		Name: message.MessageName(),
		Data: message.Data(),
	})
	if err != nil {
		return err
	}
	return s.nc.Publish(topicName, data)
}

func (s *Stream) Subscribe(topicName string, handler am.RawMessageHandler) (err error) {
	_, err = s.nc.Subscribe(topicName, s.handleMsg(topicName, handler))

	return
}

func (s *Stream) handleMsg(topicName string, handler am.RawMessageHandler) func(*nats.Msg) {
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
		}

		wCtx, cancel := context.WithCancel(context.Background())
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