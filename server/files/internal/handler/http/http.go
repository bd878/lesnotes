package http

import (
	"io"
	"net/http"
	"context"

	files "github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	SaveFileStream(ctx context.Context, stream io.Reader, id, userID int64, fileName, description string, private bool, mime string) (err error)
	ReadFileStream(ctx context.Context, id int64, fileName string, public bool) (meta *files.File, reader io.Reader, err error)
	ListFiles(ctx context.Context, userID int64, limit, offset int32, ascending, private bool) (list *files.List, err error)
	PublishFile(ctx context.Context, id, userID int64) (err error)
	PrivateFile(ctx context.Context, id, userID int64) (err error)
	DeleteFile(ctx context.Context, id, userID int64) (err error)
	ReadFileMeta(ctx context.Context, id, userID int64, public bool) (file *files.File, err error)
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
