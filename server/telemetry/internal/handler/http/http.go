package http

import (
	"io"
	"net/http"

	"github.com/bd878/gallery/server/logger"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) SendLog(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	defer req.Body.Close()

	logger.Error(data)
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	io.WriteString(w, "ok\n")
	return nil
}
