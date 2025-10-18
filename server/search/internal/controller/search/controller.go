package search

import (
	"context"

	"github.com/bd878/gallery/server/logger"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

type MessagesRepository interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) error
	DeleteMessage(ctx context.Context, id, userID int64) error
	SearchMessages(ctx context.Context, userID int64, substr string) (list []*searchmodel.Message, err error)
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

func (c *Controller) SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) (err error) {
	logger.Debugw("save search message", "id", id, "user_id", userID, "name", name, "title", title, "text", text, "private", private)

	return c.messagesRepo.SaveMessage(ctx, id, userID, name, title, text, private)
}

func (c *Controller) DeleteMessage(ctx context.Context, id, userID int64) (err error) {
	logger.Debugw("delete search message", "id", id, "user_id", userID)

	return c.messagesRepo.DeleteMessage(ctx, id, userID)
}

func (c *Controller) SearchMessages(ctx context.Context, userID int64, substr string) (list []*searchmodel.Message, err error) {
	logger.Debugw("search messages", "user_id", userID, "substr", substr)

	return c.messagesRepo.SearchMessages(ctx, userID, substr)
}