package http

import (
	"io"
	"errors"
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

func (h *Handler) UploadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	if err := req.ParseMultipartForm(1); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse form",
		})

		return err
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return errors.New("user required")
	}

	if _, ok := req.MultipartForm.File["file"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "file required",
		})

		return errors.New("file required")
	}

	f, fh, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot read file",
		})

		return errors.New("cannot read file")
	}

	fileName := filepath.Base(fh.Filename)

	fileResult, err := h.controller.SaveFileStream(req.Context(), log, f, &model.SaveFileParams{
		UserID: user.ID,
		Name:   fileName,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot save file",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.UploadFileServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "saved",
		},
		ID: fileResult.ID,
		Name: fileName,
	})

	return nil
}

func (h *Handler) DownloadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
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
			return err
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
		return errors.New("user required")
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), log, &model.ReadFileStreamParams{FileID: fileID, UserID: user.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to read file",
		})
		return err
	}

	log.Infow("downloading file", "name", file.Name)

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

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	io.WriteString(w, "ok\n")
	return nil
}
