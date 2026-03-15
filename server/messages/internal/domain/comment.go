package domain

import (
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	CommentCreatedEvent = "messages.CommentCreated"
	CommentDeletedEvent = "messages.CommentDeleted"
	CommentUpdatedEvent = "messages.CommentUpdated"
	MessageCommentsDeletedEvent = "messages.MessageCommentsDeleted"
)

type CommentCreated struct {
	ID        int64
	UserID    int64
	MessageID int64
	Text      string
	CreatedAt string
	UpdatedAt string
}

func (CommentCreated) Key() string { return CommentCreatedEvent }

func CreateComment(id, userID, messageID int64, text string, createdAt, updatedAt string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(CommentCreatedEvent, &CommentCreated{
		ID:           id,
		UserID:       userID,
		MessageID:    messageID,
		Text:         text,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}), nil
}

type CommentDeleted struct {
	ID      int64
	UserID  int64
}

func (CommentDeleted) Key() string { return CommentDeletedEvent }

func DeleteComment(id, userID int64) (ddd.Event, error) {
	return ddd.NewEvent(CommentDeletedEvent, &CommentDeleted{
		ID:      id,
		UserID:  userID,
	}), nil
}

type CommentUpdated struct {
	ID         int64
	UserID     int64
	Text       *string
	UpdatedAt  string
}

func (CommentUpdated) Key() string { return CommentUpdatedEvent }

func UpdateComment(id, userID int64, text *string, updatedAt string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/*TODO: other errors*/

	return ddd.NewEvent(CommentUpdatedEvent, &CommentUpdated{
		ID:        id,
		UserID:    userID,
		Text:      text,
		UpdatedAt: updatedAt,
	}), nil
}

type MessageCommentsDeleted struct {
	MessageID int64
}

func (MessageCommentsDeleted) Key() string { return MessageCommentsDeletedEvent }

func DeleteMessageComments(messageID int64) (ddd.Event, error) {
	return ddd.NewEvent(MessageCommentsDeletedEvent, &MessageCommentsDeleted{
		MessageID:     messageID,
	}), nil
}
