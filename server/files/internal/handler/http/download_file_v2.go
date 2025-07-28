package http

import (
	"io"
	"errors"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DownloadFileV2(w http.ResponseWriter, req *http.Request) error {
	userIDStr, fileName := req.PathValue("user_id"), req.PathValue("name")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user_id required",
		})
		return errors.New("no user_id in path given")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "error parsing request",
		})
		return err
	}

	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "name required",
		})
		return errors.New("no name in path given")
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), &model.ReadFileStreamParams{FileName: fileName, UserID: int32(userID), Public: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to read file",
		})
		return err
	}

	logger.Infow("downloading file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to write file to response",
		})
		return err
	}

	return nil
}