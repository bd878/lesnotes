package http

import (
	"net/http"
	"io"
	"context"
	"encoding/json"

	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/logger"
)

type Controller interface {
	ReadOneMessage(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (*model.Message, error)
	SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) (*model.SaveMessageResult, error)
	UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (*model.DeleteMessageResult, error)
	DeleteMessages(ctx context.Context, log *logger.Logger, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error)
	PublishMessages(ctx context.Context, log *logger.Logger, params *model.PublishMessagesParams) (*model.PublishMessagesResult, error)
	PrivateMessages(ctx context.Context, log *logger.Logger, params *model.PrivateMessagesParams) (*model.PrivateMessagesResult, error)
	ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (*model.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
}

type FilesGateway interface {
	ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (*model.ReadBatchFilesResult, error)
	ReadFile(ctx context.Context, log *logger.Logger, userID, fileID int32) (*filesmodel.File, error)
	SaveFile(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
}

type Handler struct {
	controller    Controller
	filesGateway  FilesGateway
}

func New(controller Controller, filesGateway FilesGateway) *Handler {
	return &Handler{controller, filesGateway}
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
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
