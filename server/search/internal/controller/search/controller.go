package search

import (
	"context"

	"github.com/bd878/gallery/server/logger"
)

type MessagesRepository interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string) error
	DeleteMessage(ctx context.Context, id, userID int64) error
}

type FilesRepository interface {}

type Config struct {}

type Controller struct {
	conf         Config
	messagesRepo MessagesRepository
	filesRepo    FilesRepository
}

func New(conf Config, messagesRepo MessagesRepository, filesRepo FilesRepository) *Controller {
	return &Controller{conf, messagesRepo, filesRepo}
}

func (c *Controller) SaveMessage(ctx context.Context, id, userID int64, name, title, text string) (err error) {
	logger.Debugw("save search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text)

	return c.messagesRepo.SaveMessage(ctx, id, userID, name, title, text)
}

func (c *Controller) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("delete search message", "id", id, "user_id", userID)

	return c.messagesRepo.DeleteMessage(ctx, id, userID)
}