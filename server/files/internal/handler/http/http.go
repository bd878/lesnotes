package http

import (
	"io"
	"fmt"
	"strconv"
	"net/http"
	"context"
	"path/filepath"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

type Controller interface {
	SaveFileStream(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
	ReadFileStream(ctx context.Context, log *logger.Logger, params *model.ReadFileStreamParams) (*model.File, io.Reader, error)
}

type Handler struct {
	controller  Controller
}

func New(controller Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) UploadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(1); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	if _, ok := req.MultipartForm.File["file"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "file required",
		})
		return
	}

	f, fh, err := req.FormFile("file")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot read file",
		})
		return
	}

	fileName := filepath.Base(fh.Filename)

	fileResult, err := h.controller.SaveFileStream(req.Context(), log, f, &model.SaveFileParams{
		UserID: user.ID,
		Name:   fileName,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot save file",
		})
		return
	}

	json.NewEncoder(w).Encode(model.UploadFileServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "saved",
		},
		ID: fileResult.ID,
		Name: fileName,
	})
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
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), log, &model.ReadFileStreamParams{FileID: fileID, UserID: user.ID})
	if err != nil {
		log.Errorw("failed to read file stream", "id", fileID, "user_id", user.ID, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to write file to response",
		})
		return
	}
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok\n")
}
