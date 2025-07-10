package http

import (
	"net/http"

	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) UploadFileV2(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	return h.uploadFile(log, w, req)
}