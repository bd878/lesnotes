package http

import (
	"net/http"
	"strconv"
	"path/filepath"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) SendMessage(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		threadID int64
		fileIDs  []int64
		private, hasFile bool
	)

	if err = req.ParseMultipartForm(50 << 20) /* 50 MB */; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:     server.CodeNoForm,
				Explain: "failed to parse form",
			},
		})

		return err
	}

	user, ok := req.Context().Value(middleware.UserContextKey{}).(*users.User)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "user required",
			},
		})

		return
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("file_ids") != "" {
		fileIDs = make([]int64, 0)

		if err = json.Unmarshal([]byte(req.PostFormValue("file_ids")), &fileIDs); err != nil {		
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    messages.CodeWrongFileID,
					Explain: "invalid file_ids",
				},
			})

			return
		}
	}

	if req.PostFormValue("thread") != "" {
		id, err := strconv.Atoi(req.PostFormValue("thread"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    messages.CodeWrongThreadID,
					Explain: "invalid thread",
				},
			})

			return err
		}

		threadID = int64(id)
	}

	if req.PostFormValue("public") != "" {
		public, err := strconv.Atoi(req.PostFormValue("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    messages.CodeWrongPublic,
					Explain: "invalid public",
				},
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

	if user.ID == users.PublicUserID {
		private = false
	}

	if _, ok := req.MultipartForm.File["file"]; ok {
		hasFile = true
	}

	if fileIDs == nil && !hasFile && text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:     server.CodeWrongFormat,
				Explain: "text or file_id or file required",
			},
		})

		return
	}

	if fileIDs != nil && hasFile {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "should be either file_id or file, not both",
			},
		})

		return nil
	}

	// TODO: move file saving logic in controller
	if hasFile {
		f, fh, err := req.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: messages.CodeNoFile,
					Explain: "cannot read file",
				},
			})

			return err
		}

		fileName := filepath.Base(fh.Filename)

		fileResult, err := h.filesGateway.SaveFile(req.Context(), f, &messages.SaveFileParams{
			UserID: user.ID,
			Name:   fileName,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    messages.CodeSaveFileFailed,
					Explain: "cannot save file",
				},
			})

			return err
		}

		fileIDs = append(fileIDs, int64(fileResult.ID))
	}

	message := &messages.Message{
		Text:     text,
		FileIDs:  fileIDs,
		ThreadID: threadID,
		UserID:   user.ID,
		Private:  private,
	}

	return h.saveMessage(w, req, message)
}

func (h *Handler) saveMessage(w http.ResponseWriter, req *http.Request, message *messages.Message) error {
	resp, err := h.controller.SaveMessage(req.Context(), message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: messages.CodeSaveFailed,
				Explain: "failed to save message",
			},
		})

		return err
	}

	message.ID = resp.ID
	message.Private = resp.Private
	message.UpdateUTCNano = resp.UpdateUTCNano
	message.CreateUTCNano = resp.CreateUTCNano

	var list []*files.File
	for _, id := range message.FileIDs {
		file, err := h.filesGateway.ReadFile(req.Context(), message.UserID, id)
		if err != nil {
			logger.Errorw("failed to read file for a message", "user_id", message.UserID, "file_id", id, "message_id", message.ID)
			continue
		}

		list = append(list, &files.File{
			Name: file.Name,
			ID:   file.ID,
		})
	}
	message.Files = list

	// TODO: load a user to message by UserID

	response, err := json.Marshal(messages.SendResponse{
		Message:   message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return nil
}