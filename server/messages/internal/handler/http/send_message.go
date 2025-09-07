package http

import (
	"io"
	"bytes"
	"net/http"
	"strconv"
	"path/filepath"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
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
	title := req.PostFormValue("title")

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

	if fileIDs == nil && !hasFile && text == "" && title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:     server.CodeWrongFormat,
				Explain: "text or file_id or file or title required",
			},
		})

		return
	}

	if fileIDs != nil && hasFile {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
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
					Code:    messages.CodeNoFile,
					Explain: "cannot read file",
				},
			})

			return err
		}

		fileName := filepath.Base(fh.Filename)

		var buf bytes.Buffer
		io.CopyN(&buf, f, 512)
		mime := http.DetectContentType(buf.Bytes())
		f.Seek(0, io.SeekStart)

		id := int64(utils.RandomID())

		err = h.filesGateway.SaveFile(req.Context(), f, id, user.ID, fileName, private, mime)
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

		fileIDs = append(fileIDs, id)
	}

	id := utils.RandomID()
	name := uuid.New().String()

	return h.saveMessage(w, req, int64(id), text, title, fileIDs, threadID, user.ID, private, name)
}

func (h *Handler) saveMessage(w http.ResponseWriter, req *http.Request, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) error {
	message, err := h.controller.SaveMessage(req.Context(), id, text, title, fileIDs, threadID, userID, private, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodeSaveFailed,
				Explain: "failed to save message",
			},
		})

		return err
	}

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