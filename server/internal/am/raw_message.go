package am

import (
	"time"

	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	RawMessageHandler = MessageHandler

	RawMessageSubscriber = MessageSubscriber

	RawMessageStream = MessageStream

	RawMessage struct {
		id    string
		name  string
		data  []byte
		metadata ddd.Metadata
		sentAt time.Time
		subject string
	}
)

var _ Message = (*RawMessage)(nil)

func NewRawMessage(id, name string, data []byte) *RawMessage {
	return &RawMessage{id: id, name: name, data: data, metadata: make(ddd.Metadata)}
}

func (m RawMessage) ID() string { return m.id }
func (m RawMessage) MessageName() string { return m.name }
func (m RawMessage) Data() []byte { return m.data }
func (m RawMessage) Metadata() ddd.Metadata { return m.metadata }
func (m RawMessage) SentAt() time.Time { return m.sentAt }
func (m RawMessage) Subject() string { return m.subject }
