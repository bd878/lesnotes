package http

import (
	"io"
	"errors"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DownloadFileV2(w http.ResponseWriter, req *http.Request) (err error) {
	userIDStr, fileName := req.PathValue("user_id"), req.PathValue("name")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user_id required",
			},
		})
		return errors.New("no user_id in path given")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "error parsing request",
			},
		})

		return err
	}

	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: files.CodeNoFileName,
				Explain: "name required",
			},
		})
		return errors.New("no name in path given")
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), &files.ReadFileStreamParams{FileName: fileName, UserID: int64(userID), Public: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error:  &server.ErrorCode{
				Code:  files.CodeReadFailed,
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
			Error:  &server.ErrorCode{
				Code:    files.CodeWriteFailed,
				Explain: "failed to write file to response",
			},
		})
		return err
	}

	return nil
}