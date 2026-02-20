package http

import (
	"io"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/internal/logger"
	server "github.com/bd878/gallery/server/pkg/model"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) SendLog(w http.ResponseWriter, req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:   server.CodeReadFailed,
				Explain: "cannot read",
			},
		})
		return err
	}

	defer req.Body.Close()

	logger.Infoln(string(data))
	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
	})
	return nil
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) error {
	io.WriteString(w, "ok\n")
	return nil
}
