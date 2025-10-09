package http

import (
	"io"
	"net/http"
)

type Controller interface {
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
