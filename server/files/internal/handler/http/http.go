package http

import (
	"io"
	"fmt"
	"strconv"
	"net/http"
	"context"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/files/pkg/model"
)

type Controller interface {
	ReadFileStream(ctx context.Context, log *logger.Logger, params *model.ReadFileStreamParams) (*model.File, io.Reader, error)
}

type Handler struct {
	controller  Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) DownloadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	var (
		fileID int32
	)

	values := req.URL.Query()
	if values.Get("id") != "" {
		fileid, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ServerResponse{
				Status: "ok",
				Description: "id is empty",
			})
			return
		}
		fileID = int32(fileid)
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), log, &model.ReadFileStreamParams{FileID: fileID, UserID: user.ID})
	if err != nil {
		log.Errorw("failed to read file stream", "id", fileID, "user_id", user.ID, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status:      "error",
			Description: "failed to read file",
		})
		return
	}

	log.Infow("downloading file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		log.Errorw("failed to write file stream to response", "id", fileID, "user_id", user.ID, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status:      "error",
			Description: "failed to write file to response",
		})
		return
	}
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok\n")
}
