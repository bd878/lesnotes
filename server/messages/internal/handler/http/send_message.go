package http

import (
	"net/http"
	"strconv"
	"fmt"
	"path/filepath"
	"encoding/json"

	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var (
		err error
		fileID, threadID int32
		private, hasFile bool
	)

	if err = req.ParseMultipartForm(50 << 10) /* 50 MB */; err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

		return fmt.Errorf("no user")
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("file_id") != "" {
		fileid, err := strconv.Atoi(req.PostFormValue("file_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid file_id",
			})

			return err
		}

		fileID = int32(fileid)
	}

	values := req.URL.Query()
	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid thread id",
			})

			return err
		}

		threadID = int32(threadid)
	}

	if values.Has("public") {
		public, err := strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid public param",
			})

			return err
		}

		if public > 0 {
			private = false
		} else if public == 0 {
			private = true
		} else {
			private = true
		}
	} else {
		private = true
	}

	if user.ID == usermodel.PublicUserID {
		private = false
	}

	if _, ok := req.MultipartForm.File["file"]; ok {
		hasFile = true
	}

	log.Infow("received message", "text", text, "name", user.Name, "file_id", fileID, "has_file", hasFile)

	if fileID == 0 && !hasFile && text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "text or file_id or file required",
		})

		return nil
	}

	if fileID != 0 && hasFile {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "should be either file_id or file, not both",
		})

		return nil
	}

	if hasFile {
		f, fh, err := req.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot read file",
			})

			return err
		}

		fileName := filepath.Base(fh.Filename)

		fileResult, err := h.filesGateway.SaveFile(req.Context(), log, f, &model.SaveFileParams{
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

		fileID = fileResult.ID
	}

	message := &model.Message{
		Text: text,
		FileID: fileID,
		ThreadID: threadID,
		UserID: user.ID,
		Private: private,
	}

	return h.saveMessage(log, w, req, message)
}

func (h *Handler) saveMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request, message *model.Message) error {
	resp, err := h.controller.SaveMessage(req.Context(), log, message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to save message",
		})

		return err
	}

	message.ID = resp.ID
	message.Private = resp.Private
	message.UpdateUTCNano = resp.UpdateUTCNano
	message.CreateUTCNano = resp.CreateUTCNano

	if message.FileID != 0 {
		fileRes, err := h.filesGateway.ReadFile(req.Context(), log, message.UserID, message.FileID)
		if err != nil {
			log.Errorw("failed to read file for a message", "user_id", message.UserID, "file_id", message.FileID, "message_id", resp.ID)
		} else {
			message.File = &filesmodel.File{
				Name: fileRes.Name,
				ID: message.FileID,
			}
		}
	}

	json.NewEncoder(w).Encode(model.NewMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		Message: message,
	})

	return nil
}