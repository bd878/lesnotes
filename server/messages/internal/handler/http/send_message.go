package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	"github.com/bd878/gallery/server/internal/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) SendMessage(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		threadID         int64
		fileIDs          []int64
		private          bool
	)

	if err = req.ParseMultipartForm(50 << 20); /* 50 MB */ err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoForm,
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
				Error: &server.ErrorCode{
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

	if fileIDs == nil && text == "" && title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "text or file_id or file or title required",
			},
		})

		return
	}

	id := utils.RandomID()
	name := utils.RandomString(8)

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

	// TODO: load a user to message by UserID

	response, err := json.Marshal(messages.SendResponse{
		Message: message,
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
