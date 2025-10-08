package am

type (
	RawMessageStream = MessageStream[RawMessage, IncomingRawMessage]

	RawMessage interface {
		Message
		Data() []byte
	}

	IncomingRawMessage interface {
		IncomingMessage
		Data() []byte
	}

	RawMessageHandler = MessageHandler[IncomingRawMessage]

	rawMessage struct {
		id      string
		name    string
		subject string
		data    []byte
	}
)

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string { return m.id }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte { return m.data }
func (m rawMessage) Subject() string { return m.subject }
