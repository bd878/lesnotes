package domain

const (
	MessageCreatedEvent = "messages.MessageCreated"
	MessageDeletedEvent = "messages.MessageDeleted"
)

type MessageCreated struct {
	ID        int64
	Text      string
	Title     string
	Name      string
}

func (MessageCreated) Key() string { return MessageCreatedEvent }

type MessageDeleted struct {
	ID        int64
}

func (MessageDeleted) Key() string { return MessageDeletedEvent }
