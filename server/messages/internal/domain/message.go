package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const MessageAggregate = "messages.Message"

var (
	ErrIDRequired = errors.New("id is 0")
)

func CreateMessage(id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(MessageCreatedEvent, &MessageCreated{
		ID:     id,
		Text:   text,
		Title:  title,
		Name:   name,
	}), nil
}

func DeleteMessage(id int64) (ddd.Event, error) {
	return ddd.NewEvent(MessageDeletedEvent, &MessageDeleted{
		ID:     id,
	}), nil
}