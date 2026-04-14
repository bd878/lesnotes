package domain

import (
	"errors"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	FileUploadedEvent    = "files.FileUploaded"
	FilesDeletedEvent    = "files.FilesDeleted"
	FilesPublishedEvent  = "files.FilesPublished"
	FilesPrivatedEvent   = "files.FilesPrivated"
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
	CreatedAt    string
	UpdatedAt    string
}

func (FileUploaded) Key() string { return FileUploadedEvent }

func UploadFile(id int64, name, description string, userID int64, private bool, mime string, size int64, createdAt, updatedAt string) (ddd.Event, error) {
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
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}), nil
}

type FilesDeleted struct {
	IDs      []int64
	UserID   int64
}

func (FilesDeleted) Key() string { return FilesDeletedEvent }

func DeleteFiles(userID int64, ids []int64) (ddd.Event, error) {
	if len(ids) == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(FilesDeletedEvent, &FilesDeleted{
		IDs:       ids,
		UserID:    userID,
	}), nil
}

type FilesPublished struct {
	IDs        []int64
	UserID     int64
	UpdatedAt  string
}

func (FilesPublished) Key() string { return FilesPublishedEvent }

func PublishFiles(userID int64, ids []int64, updatedAt string) (ddd.Event, error) {
	if len(ids) == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(FilesPublishedEvent, &FilesPublished{IDs: ids, UserID: userID, UpdatedAt: updatedAt}), nil
}

type FilesPrivated struct {
	IDs        []int64
	UserID     int64
	UpdatedAt  string
}

func (FilesPrivated) Key() string { return FilesPrivatedEvent}

func PrivateFiles(userID int64, ids []int64, updatedAt string) (ddd.Event, error) {
	if len(ids) == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(FilesPrivatedEvent, &FilesPrivated{IDs: ids, UserID: userID, UpdatedAt: updatedAt}), nil
}
