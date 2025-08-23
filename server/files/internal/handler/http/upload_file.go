package http

import (
	"errors"
	"net/http"
	"encoding/json"
	"path/filepath"

	"github.com/bd878/gallery/server/utils"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) UploadFile(w http.ResponseWriter, req *http.Request) error {
	return h.uploadFile(w, req, 0)
}

func (h *Handler) uploadFile(w http.ResponseWriter, req *http.Request, public int) (err error) {
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
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user required",
			},
		})

		return errors.New("user required")
	}

	if user.ID == users.PublicUserID {
		private = false
	}

	if err := req.ParseMultipartForm(50 << 20) /* 50 MB */; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "failed to parse form",
			},
		})

		return err
	}

	if _, ok := req.MultipartForm.File["file"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: files.CodeNoFile,
				Explain: "file required",
			},
		})

		return errors.New("file required")
	}

	f, fh, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: files.CodeReadFailed,
				Explain: "cannot read file",
			},
		})

		return errors.New("cannot read file")
	}

	fileName := filepath.Base(fh.Filename)

	fileResult, err := h.controller.SaveFileStream(req.Context(), f, &files.SaveFileParams{
		UserID: user.ID,
		Name:   fileName,
		Private: private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: files.CodeSaveFailed,
				Explain: "cannot save file",
			},
		})

		return err
	}

	response, err := json.Marshal(files.UploadResponse{
		ID: fileResult.ID,
		Name: fileName,
		Description: "saved",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return nil
}