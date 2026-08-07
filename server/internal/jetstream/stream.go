package jetstream

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/structpb"
	"github.com/nats-io/nats.go"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/am"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
	mu         sync.Mutex
	subs       []*nats.Subscription
}

var _ am.MessageStream = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext) *Stream {
	return &Stream{
		streamName: streamName,
		js:         js,
	}
}

func (s *Stream) Publish(ctx context.Context, topicName string, rawMsg am.Message) (err error) {
	var data []byte

	metadata, err := structpb.NewStruct(rawMsg.Metadata())
	if err != nil {
		return err
	}

	data, err = proto.Marshal(&api.StreamMessage{
		Id:      rawMsg.ID(),
		Name:    rawMsg.MessageName(),
		Data:    rawMsg.Data(),
		Metadata: metadata,
		SentAt: timestamppb.New(rawMsg.SentAt()),
	})
	if err != nil {
		return
	}

	slog.Debug("subject", slog.String("subject", rawMsg.Subject()))

	var p nats.PubAckFuture
	p, err = s.js.PublishMsgAsync(&nats.Msg{
		Subject: rawMsg.Subject(),
		Data:    data,
	}, nats.MsgId(rawMsg.ID()))
	if err != nil {
		return
	}

	go func(future nats.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok():
				return
			case <-future.Err():
				tries = tries - 1
				if tries <= 0 {
					slog.Error("unable to publish message", slog.Int("maxRetries", maxRetries))
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					slog.Error("publish async retry failed", slog.String("error", err.Error()))
					return
				}
			}
		}
	}(p, maxRetries)

	return
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler, options ...am.SubscriberOption) (am.Subscription, error) {
	var err error

	s.mu.Lock()
	defer s.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}
	cfg := &nats.ConsumerConfig{
		MaxDeliver: subCfg.MaxRedeliver(),
		DeliverSubject: topicName,
		FilterSubject: topicName,
	}
	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.DeliverSubject = groupName
		cfg.DeliverGroup = groupName
		cfg.Durable = groupName

		opts = append(opts, nats.Bind(s.streamName, groupName), nats.Durable(groupName))
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		ackWait := subCfg.AckWait()

		cfg.AckPolicy = nats.AckExplicitPolicy
		cfg.AckWait = ackWait

		opts = append(opts, nats.AckExplicit(), nats.AckWait(ackWait))
	} else {
		cfg.AckPolicy = nats.AckNonePolicy
		opts = append(opts, nats.AckNone())
	}

	_, err = s.js.AddConsumer(s.streamName, cfg)
	if err != nil {
		return nil, err
	}

	var sub *nats.Subscription

	if groupName := subCfg.GroupName(); groupName == "" {
		sub, err = s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
	} else {
		sub, err = s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
	}

	s.subs = append(s.subs, sub)

	return subscription{sub}, nil
}

func (s *Stream) Unsubscribe() error {
	for _, sub := range s.subs {
		if !sub.IsValid() {
			continue
		}
		err := sub.Drain()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler) func(*nats.Msg) {
	return func(natsMsg *nats.Msg) {
		var err error

		m := &api.StreamMessage{}
		err = proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			slog.Warn("unmarshal error", slog.String("error", err.Error()))
			return
		}

		msg := &rawMessage{
			id:    m.GetId(),
			name:  m.GetName(),
			subject: natsMsg.Subject,
			data: m.GetData(),
			metadata: m.GetMetadata().AsMap(),
			sentAt: m.SentAt.AsTime(),
			receivedAt: time.Now(),
			acked: false,
			ackFn: func() error { return natsMsg.Ack() },
			nackFn: func() error { return natsMsg.Nak() },
			extendFn: func() error { return natsMsg.InProgress() },
			killFn: func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				slog.Warn("auto ack error", slog.String("error", err.Error()))
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					slog.Warn("ack error", slog.String("error", ackErr.Error()))
				}
				return
			}
			slog.Error("handler error", slog.String("error", err.Error()))
			if nakErr := msg.NAck(); nakErr != nil {
				slog.Warn("nack error", slog.String("error", nakErr.Error()))
			}
		case <-wCtx.Done():
			return
		}
	}
}