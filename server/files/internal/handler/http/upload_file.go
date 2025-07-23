package http

import (
	"errors"
	"net/http"
	"encoding/json"
	"path/filepath"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) UploadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	return h.uploadFile(log, w, req, 0)
}

func (h *Handler) uploadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request, public int) error {
	var private bool

	if public > 0 {
		private = false
	} else if public == 0 {
		private = true
	} else {
		private = true
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

	if user.ID == usermodel.PublicUserID {
		private = false
	}

	if err := req.ParseMultipartForm(50 << 10) /* 50 MB */; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse form",
		})

		return err
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
		Private: private,
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