package http

import (
	"io"
	"context"
	"net/http"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
)

type Controller interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []int64, isLastPage bool, err error)
	ReadThread(ctx context.Context, id, userID int64) (thread *threads.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	PublishThread(ctx context.Context, id, userID int64) (err error)
	PrivateThread(ctx context.Context, id, userID int64) (err error)
	CreateThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name string, private bool) (err error)
	DeleteThread(ctx context.Context, id, userID int64) (err error)
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error)
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
