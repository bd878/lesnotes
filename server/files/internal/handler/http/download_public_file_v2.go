package http

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) DownloadPublicFileV2(w http.ResponseWriter, req *http.Request) (err error) {
	fileName := req.PathValue("name")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    files.CodeNoFileName,
				Explain: "name required",
			},
		})

		return
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), 0, users.PublicUserID, fileName, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error:  &server.ErrorCode{
				Code:    files.CodeReadFailed,
				Explain: "failed to read file",
			},
		})
		return err
	}

	logger.Infow("downloading public file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error:  &server.ErrorCode{
				Code:    files.CodeWriteFailed,
				Explain: "failed to write file to response",
			},
		})
		return
	}

	return
}