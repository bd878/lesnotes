package domain

const (
	MessageCreatedEvent = "messages.MessageCreated"
	MessageDeletedEvent = "messages.MessageDeleted"
)

type MessageCreated struct {
	Message *Message
}

func (MessageCreated) Key() string { return MessageCreatedEvent }

type MessageDeleted struct {
	Message *Message
}

func (MessageDeleted) Key() string { return MessageDeletedEvent }
