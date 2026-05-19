package nats

import (
	"time"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	rawMessage struct {
		id      string
		data    []byte
		name    string
		subject string
		metadata ddd.Metadata
		sentAt time.Time
	}
)

var _ am.IncomingMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string { return m.id }
func (m rawMessage) Data() []byte { return m.data }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Subject() string { return m.subject }
func (m rawMessage) Metadata() ddd.Metadata { return m.metadata }
func (m rawMessage) SentAt() time.Time { return m.sentAt }
