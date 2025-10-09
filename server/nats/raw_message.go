package nats

import (
	"github.com/bd878/gallery/server/am"
)

type (
	rawMessage struct {
		id      string
		data    []byte
		name    string
		subject string
		ackFn   func() error
		nakFn   func() error
	}
)

var _ am.IncomingMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string { return m.id }
func (m rawMessage) Data() []byte { return m.data }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Subject() string { return m.subject }

func (m *rawMessage) Ack() error { return m.ackFn() }
func (m *rawMessage) Nak() error { return m.nakFn() }
