package http

import (
	"net/http"
	"io"
	"context"
	"encoding/json"

	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type Controller interface {
	ReadOneMessage(ctx context.Context, params *model.ReadOneMessageParams) (*model.Message, error)
	SaveMessage(ctx context.Context, message *model.Message) (*model.SaveMessageResult, error)
	UpdateMessage(ctx context.Context, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, params *model.DeleteMessageParams) (*model.DeleteMessageResult, error)
	DeleteMessages(ctx context.Context, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error)
	PublishMessages(ctx context.Context, params *model.PublishMessagesParams) (*model.PublishMessagesResult, error)
	PrivateMessages(ctx context.Context, params *model.PrivateMessagesParams) (*model.PrivateMessagesResult, error)
	ReadAllMessages(ctx context.Context, params *model.ReadMessagesParams) (*model.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
}

type FilesGateway interface {
	ReadBatchFiles(ctx context.Context, params *model.ReadBatchFilesParams) (*model.ReadBatchFilesResult, error)
	ReadFile(ctx context.Context, userID, fileID int64) (*filesmodel.File, error)
	SaveFile(ctx context.Context, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to get status",
		})

		return err
	}

	return nil
}
