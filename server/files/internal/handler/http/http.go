package http

import (
	"io"
	"net/http"
	"context"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	SaveFileStream(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
	ReadFileStream(ctx context.Context, log *logger.Logger, params *model.ReadFileStreamParams) (*model.File, io.Reader, error)
}

type Handler struct {
	controller  Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	io.WriteString(w, "ok\n")
	return nil
}
