package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const (
	MessageCreatedEvent = "messages.MessageCreated"
	MessageDeletedEvent = "messages.MessageDeleted"
	MessageUpdatedEvent = "messages.MessageUpdated"
	MessagePublishEvent = "messages.MessagePublished"
	MessagePrivateEvent = "messages.MessagePrivated"
)

var (
	ErrIDRequired = errors.New("id is 0")
)

type MessageCreated struct {
	ID        int64
	UserID    int64
	Text      string
	Title     string
	Name      string
	Private   bool
}

func (MessageCreated) Key() string { return MessageCreatedEvent }

func CreateMessage(id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(MessageCreatedEvent, &MessageCreated{
		ID:      id,
		UserID:  userID,
		Text:    text,
		Title:   title,
		Name:    name,
		Private: private,
	}), nil
}

type MessageDeleted struct {
	ID        int64
	UserID    int64
	Name      string
}

func (MessageDeleted) Key() string { return MessageDeletedEvent }

func DeleteMessage(id, userID int64) (ddd.Event, error) {
	return ddd.NewEvent(MessageDeletedEvent, &MessageDeleted{
		ID:     id,
		UserID: userID,
	}), nil
}
