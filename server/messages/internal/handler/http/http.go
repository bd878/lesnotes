package http

import (
	"io"
	"context"
	"net/http"

	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *messages.Message, err error)
	UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, threadID int64, userID int64, private int32) (err error)
	DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PublishMessages(ctx context.Context, ids []int64, userID int64) (err error)
	PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error)
	ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *messages.Message, err error)
	ReadMessagesAround(ctx context.Context, userID, threadID, id int64, limit int32) (list *messages.List, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (list *messages.List, err error)
	ReadThreadMessages(ctx context.Context, userID int64, threadID int64, limit, offset int32, ascending bool) (list *messages.List, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*messages.Message, err error)
	ReadPath(ctx context.Context, userID int64, id int64) (messages []*messages.Message, err error)
}

// TODO: think to move on controller/service level
type FilesGateway interface {
	ReadBatchFiles(ctx context.Context, fileIDs []int64, userID int64) (files map[int64]*files.File, err error)
	ReadFile(ctx context.Context, userID, fileID int64) (file *files.File, err error)
	SaveFile(ctx context.Context, stream io.Reader, id, userID int64, name string, private bool, mime string) (err error)
}

type Handler struct {
	controller    Controller
	filesGateway  FilesGateway
}

func New(controller Controller, filesGateway FilesGateway) *Handler {
	return &Handler{controller, filesGateway}
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
