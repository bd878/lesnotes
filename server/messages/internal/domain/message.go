package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const (
	MessageCreatedEvent = "messages.MessageCreated"
	MessageDeletedEvent = "messages.MessageDeleted"
	// TODO: user deleted: delete all messages
	MessageUpdatedEvent = "messages.MessageUpdated"
	MessagesPublishEvent = "messages.MessagesPublished"
	MessagesPrivateEvent = "messages.MessagesPrivated"
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

type MessageUpdated struct {
	ID        int64
	UserID    int64
	Text      string
	Title     string
	Name      string
	Private   int32
}

func (MessageUpdated) Key() string { return MessageUpdatedEvent }

func UpdateMessage(id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private int32, name string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(MessageUpdatedEvent, &MessageUpdated{
		ID:      id,
		UserID:  userID,
		Text:    text,
		Title:   title,
		Name:    name,
		Private: private,
	}), nil
}

type MessagesPublished struct {
	IDs      []int64
	UserID   int64
}

func (MessagesPublished) Key() string { return MessagesPublishEvent }

func PublishMessages(userID int64, ids []int64) (ddd.Event, error) {
	if ids == nil {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(MessagesPublishEvent, &MessagesPublished{IDs: ids, UserID: userID}), nil
}

type MessagesPrivated struct {
	IDs      []int64
	UserID   int64
}

func (MessagesPrivated) Key() string { return MessagesPrivateEvent }

func PrivateMessages(userID int64, ids []int64) (ddd.Event, error) {
	if ids == nil {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(MessagesPrivateEvent, &MessagesPrivated{IDs: ids, UserID: userID}), nil
}