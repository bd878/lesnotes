package http

import (
	"net/http"
	"strconv"
	"fmt"
	"path/filepath"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) SendMessage(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		fileID, threadID int64
		private, hasFile bool
	)

	if err = req.ParseMultipartForm(50 << 20) /* 50 MB */; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoForm,
				Explain: "failed to parse form",
			},
		})

		return err
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

		return fmt.Errorf("no user")
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("file_id") != "" {
		fileid, err := strconv.Atoi(req.PostFormValue("file_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: messages.CodeWrongFileID,
					Explain: "invalid file_id",
				},
			})

			return err
		}

		fileID = int64(fileid)
	}

	values := req.URL.Query()
	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: messages.CodeWrongThreadID,
					Explain: "invalid thread id",
				},
			})

			return err
		}

		threadID = int64(threadid)
	}

	if values.Has("public") {
		public, err := strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code: messages.CodeWrongPublic,
					Explain: "invalid public param",
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

	if fileID == 0 && !hasFile && text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "text or file_id or file required",
			},
		})

		return nil
	}

	if fileID != 0 && hasFile {
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
					Code: messages.CodeSaveFileFailed,
					Explain: "cannot save file",
				},
			})

			return err
		}

		fileID = fileResult.ID
	}

	message := &messages.Message{
		Text: text,
		FileID: fileID,
		ThreadID: threadID,
		UserID: user.ID,
		Private: private,
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

	if message.FileID != 0 {
		fileRes, err := h.filesGateway.ReadFile(req.Context(), message.UserID, message.FileID)
		if err != nil {
			logger.Errorw("failed to read file for a message", "user_id", message.UserID, "file_id", message.FileID, "message_id", resp.ID)
		} else {
			message.File = &files.File{
				Name: fileRes.Name,
				ID: message.FileID,
			}
		}
	}

	response, err := json.Marshal(messages.SaveResponse{
		Message:   message,
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