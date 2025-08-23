package http

import (
	"io"
	"net/http"
	"context"
	"encoding/json"

	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

type Controller interface {
	SaveFileStream(ctx context.Context, stream io.Reader, params *files.SaveFileParams) (*files.SaveFileResult, error)
	ReadFileStream(ctx context.Context, params *files.ReadFileStreamParams) (*files.File, io.Reader, error)
}

type Handler struct {
	controller  Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeStatusFailed,
				Explain: "failed to get status",
			},
		})

		return err
	}

	return nil
}
