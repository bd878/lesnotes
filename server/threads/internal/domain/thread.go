package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const (
	ThreadCreatedEvent = "threads.ThreadCreated"
	ThreadDeletedEvent = "threads.ThreadDeleted"
	ThreadUpdatedEvent = "threads.ThreadUpdated"
	ThreadPublishEvent = "threads.ThreadPublished"
	ThreadPrivateEvent = "threads.ThreadPrivated"
	ThreadParentChangedEvent = "threads.ThreadParentChanged"
)

var (
	ErrIDRequired = errors.New("id is 0")
)

type ThreadCreated struct {
	ID          int64
	UserID      int64
	ParentID    int64
	Name        string
	Description string
	Private     bool
}

func (ThreadCreated) Key() string { return ThreadCreatedEvent }

func CreateThread(id, userID, parentID int64, name, description string, private bool) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	// TODO: other errors

	return ddd.NewEvent(ThreadCreatedEvent, &ThreadCreated{
		ID:          id,
		UserID:      userID,
		ParentID:    parentID,
		Name:        name,
		Description: description,
		Private:     private,
	}), nil
}

type ThreadUpdated struct {
	ID          int64
	UserID      int64
	Name        string
	Description string
}

func (ThreadUpdated) Key() string { return ThreadUpdatedEvent }

func UpdateThread(id, userID int64, name, description string) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	// TODO: other errors

	return ddd.NewEvent(ThreadUpdatedEvent, &ThreadUpdated{
		ID:             id,
		UserID:         userID,
		Name:           name,
		Description:    description,
	}), nil
}

type ThreadDeleted struct {
	ID       int64
	UserID   int64
}

func (ThreadDeleted) Key() string { return ThreadDeletedEvent }

func DeleteThread(id, userID int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(ThreadDeletedEvent, &ThreadDeleted{
		ID:      id,
		UserID:  userID,
	}), nil
}

type ThreadPublished struct {
	ID       int64
	UserID   int64
}

func (ThreadPublished) Key() string { return ThreadPublishEvent }

func PublishThread(id, userID int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(ThreadPublishEvent, &ThreadPublished{
		ID:      id,
		UserID:  userID,
	}), nil
}

type ThreadPrivated struct {
	ID         int64
	UserID     int64
}

func (ThreadPrivated) Key() string { return ThreadPrivateEvent }

func PrivateThread(id, userID int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(ThreadPrivateEvent, &ThreadPrivated{
		ID:     id,
		UserID: userID,
	}), nil
}

type ThreadParentChanged struct {
	ID          int64
	UserID      int64
	ParentID    int64
}

func (ThreadParentChanged) Key() string { return ThreadParentChangedEvent }

func ChangeThreadParent(id, userID, parentID int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(ThreadParentChangedEvent, &ThreadParentChanged{
		ID:        id,
		UserID:    userID,
		ParentID:  parentID,
	}), nil
}
