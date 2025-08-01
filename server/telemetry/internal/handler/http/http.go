package http

import (
	"io"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	servermodel "github.com/bd878/gallery/server/pkg/model"
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot read",
		})
		return err
	}

	defer req.Body.Close()

	logger.Infoln(string(data))
	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "accepted",
	})
	return nil
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	io.WriteString(w, "ok\n")
	return nil
}
