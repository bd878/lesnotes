package http

import (
	"io"
	"net/http"
	"context"

	files "github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	SaveFileStream(ctx context.Context, stream io.Reader, id, userID int64, fileName string, private bool, mime string) (err error)
	ReadFileStream(ctx context.Context, id int64, fileName string, public bool) (meta *files.File, reader io.Reader, err error)
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
		return err
	}

	return nil
}
