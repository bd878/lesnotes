package http

import (
	"io"
	"errors"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DownloadFile(w http.ResponseWriter, req *http.Request) error {
	var (
		fileID int64
	)

	values := req.URL.Query()
	if values.Get("id") != "" {
		fileid, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: server.CodeNoID,
					Explain: "id is empty",
				},
			})
			return err
		}
		fileID = int64(fileid)
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user required",
			},
		})
		return errors.New("user required")
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), &files.ReadFileStreamParams{FileID: fileID, UserID: user.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: files.CodeReadFailed,
				Explain: "failed to read file",
			},
		})
		return err
	}

	logger.Infow("downloading file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: files.CodeWriteFailed,
				Explain: "failed to write file to response",
			},
		})
		return err
	}

	return nil
}