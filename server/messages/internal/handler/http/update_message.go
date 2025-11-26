package http

import (
	"net/http"
	"context"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
)

func (h *Handler) UpdateMessage(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		id int64
		public int
		fileIDs []int64
	)

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

	if req.PostFormValue("id") != "" {
		messageID, err := strconv.Atoi(req.PostFormValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid id",
				},
			})

			return err
		}

		id = int64(messageID)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "empty message id",
			},
		})

		return
	}

	text := req.PostFormValue("text")
	title := req.PostFormValue("title")
	name := req.PostFormValue("name")

	if req.PostFormValue("file_ids") != "" {
		if err = json.Unmarshal([]byte(req.PostFormValue("file_ids")), &fileIDs); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongFormat,
					Explain: "cannot parse file_ids",
				},
			})

			return
		}
	}

	if req.PostFormValue("public") != "" {
		public, err = strconv.Atoi(req.PostFormValue("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongFormat,
					Explain: "invalid public param",
				},
			})

			return
		}
	} else {
		public = -1
	}

	return h.updateMessage(req.Context(), w, id, user, text, title, name, fileIDs, public)
}

func (h *Handler) updateMessage(ctx context.Context, w http.ResponseWriter, messageID int64, user *users.User,
	text, title, name string, fileIDs []int64, public int,
) (err error) {
	var private int32

	if public == 1 {
		private = 0
	} else if public == 0 {
		private = 1
	} else {
		private = -1
	}

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodeMessagePublic,
				Explain: "cannot move public message",
			},
		})

		return
	}

	if user.ID == users.PublicUserID && private == 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeMessagePublic,
				Explain: "cannot make private public message",
			},
		})

		return
	}

	msg, err := h.controller.ReadMessage(ctx, messageID, "", []int64{user.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodeNoMessage,
				Explain: "failed to read message",
			},
		})

		return err		
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "text must be provided",
			},
		})

		return nil
	}

	if text == "" {
		text = msg.Text
	}

	if title == "" {
		title = msg.Title
	}

	if name == "" {
		name = msg.Name
	}

	if fileIDs == nil {
		fileIDs = msg.FileIDs
	}

	err = h.controller.UpdateMessage(ctx, messageID, text, title, name, fileIDs, user.ID, private)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeUpdateFailed,
				Explain: "failed to update message",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.UpdateResponse{
		Description:   "updated",
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