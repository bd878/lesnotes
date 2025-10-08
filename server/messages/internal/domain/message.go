package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const MessageAggregate = "messages.Message"

var (
	ErrIDRequired = errors.New("id is 0")
)

type Message struct {
	ID        int64
	Text      string
	Title     string
	Name      string
}

func CreateMessage(id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(MessageCreatedEvent, &Message{
		ID:     id,
		Text:   text,
		Title:  title,
		Name:   name,
	}), nil
}