package domain

import (
	"errors"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	FileUploadedEvent   = "files.FileUploaded"
	FileDeletedEvent    = "files.FileDeleted"
	FilePublishedEvent  = "files.FilePublished"
	FilePrivatedEvent   = "files.FilePrivated"
)

var (
	ErrIDRequired = errors.New("id is 0")
)

type FileUploaded struct {
	ID           int64
	UserID       int64
	Name         string
	Description  string
	Private      bool
	Mime         string
	Size         int64
}

func (FileUploaded) Key() string { return FileUploadedEvent }

func UploadFile(id int64, name, description string, userID int64, private bool, mime string, size int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}
	/* TODO: other errors*/

	return ddd.NewEvent(FileUploadedEvent, &FileUploaded{
		ID:          id,
		UserID:      userID,
		Name:        name,
		Description: description,
		Mime:        mime,
		Size:        size,
		Private:     private,
	}), nil
}

type FileDeleted struct {
	ID      int64
	UserID  int64
}

func (FileDeleted) Key() string { return FileDeletedEvent }

func DeleteFile(id, userID int64) (ddd.Event, error) {
	return ddd.NewEvent(FileDeletedEvent, &FileDeleted{
		ID:       id,
		UserID:   userID,
	}), nil
}

type FilePublished struct {
	ID     int64
	UserID int64
}

func (FilePublished) Key() string { return FilePublishedEvent }

func PublishFile(userID, id int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(FilePublishedEvent, &FilePublished{ID: id, UserID: userID}), nil
}

type FilePrivated struct {
	ID         int64
	UserID     int64
}

func (FilePrivated) Key() string { return FilePrivatedEvent}

func PrivateFile(userID, id int64) (ddd.Event, error) {
	if id == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(FilePrivatedEvent, &FilePrivated{ID: id, UserID: userID}), nil
}
