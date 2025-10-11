package http

import (
	"io"
	"context"
	"net/http"

	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

type Controller interface {
	SearchMessages(ctx context.Context, userID int64, query string) (list []*searchmodel.Message, err error)
}

type Handler struct {
	controller Controller
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
