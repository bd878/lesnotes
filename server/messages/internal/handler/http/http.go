package http

import (
	"io"
	"context"
	"net/http"

	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	ReadOneMessage(ctx context.Context, id int64, userIDs []int64) (message *messages.Message, err error)
	SaveMessage(ctx context.Context, id int64, text string, fileIDs []int64, threadID int64, userID int64, private bool) (message *messages.Message, err error)
	UpdateMessage(ctx context.Context, params *messages.UpdateMessageParams) (*messages.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, params *messages.DeleteMessageParams) (*messages.DeleteMessageResult, error)
	DeleteMessages(ctx context.Context, params *messages.DeleteMessagesParams) (*messages.DeleteMessagesResult, error)
	PublishMessages(ctx context.Context, params *messages.PublishMessagesParams) (*messages.PublishMessagesResult, error)
	PrivateMessages(ctx context.Context, params *messages.PrivateMessagesParams) (*messages.PrivateMessagesResult, error)
	ReadAllMessages(ctx context.Context, params *messages.ReadMessagesParams) (*messages.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, params *messages.ReadThreadMessagesParams) (*messages.ReadThreadMessagesResult, error)
}

type FilesGateway interface {
	ReadBatchFiles(ctx context.Context, params *messages.ReadBatchFilesParams) (*messages.ReadBatchFilesResult, error)
	ReadFile(ctx context.Context, userID, fileID int64) (*files.File, error)
	SaveFile(ctx context.Context, stream io.Reader, params *messages.SaveFileParams) (*messages.SaveFileResult, error)
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
