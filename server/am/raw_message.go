package am

type (
	RawMessageHandler = MessageHandler[IncomingMessage]

	RawMessageSubscriber = MessageSubscriber[IncomingMessage]

	RawMessageStream = MessageStream[Message, IncomingMessage]

	RawMessage struct {
		id    string
		name  string
		data  []byte
	}
)

var _ Message = (*RawMessage)(nil)

func NewRawMessage(id, name string, data []byte) *RawMessage {
	return &RawMessage{id, name, data}
}

func (m RawMessage) ID() string { return m.id }
func (m RawMessage) MessageName() string { return m.id }
func (m RawMessage) Data() []byte { return m.data }

