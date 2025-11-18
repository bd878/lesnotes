package http

import (
	"io"
	"context"
	"net/http"
)

type Controller interface {
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
