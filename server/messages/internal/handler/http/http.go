package http

import (
	"context"
	"io"
	"net/http"

	messages "github.com/bd878/gallery/server/messages/pkg/model"
)

type MessagesController interface {
	SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *messages.Message, err error)
	UpdateMessage(ctx context.Context, id int64, text, title, name *string, fileIDs []int64, userID int64) (err error)
	DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PublishMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error)
	ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *messages.Message, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (list *messages.MessagesList, err error)
	ReadThreadMessages(ctx context.Context, userID, threadID int64, threadName string, limit, offset int32, ascending bool, privateMessage *bool) (list *messages.MessagesList, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*messages.Message, err error)
	ReadPath(ctx context.Context, userID, id int64, name string) (messages []*messages.Message, parentID int64, err error)
	ReadTree(ctx context.Context, userID, highlightID int64, highlightName string, messageID int64, name string, limit, offset int32, privateMessage *bool, pairs []*messages.IDLimitOffset) (list *messages.MessagesList, err error)
}

type TranslationsController interface {
	SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error)
	UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error)
	DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error)
	ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name *string) (translation *messages.Translation, err error)
	ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*messages.Translation, err error)
}

type CommentsController interface {
	SendComment(ctx context.Context, id, userID, messageID int64, text string, metadata []byte) (err error)
	UpdateComment(ctx context.Context, id, userID int64, text *string) (err error)
	DeleteComment(ctx context.Context, id, userID int64) (err error)
	DeleteMessageComments(ctx context.Context, messageID int64) (err error)
	ReadComment(ctx context.Context, id, userID int64) (comment *messages.Comment, err error)
	ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32, asc bool) (list *messages.CommentsList, err error)
}

type Handler struct {
	controller             MessagesController
	translationsController TranslationsController
	commentsController     CommentsController
}

func New(messagesController MessagesController, translationsController TranslationsController,
	commentsController CommentsController) *Handler {
	return &Handler{
		controller:             messagesController,
		commentsController:     commentsController,
		translationsController: translationsController,
	}
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
