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
	var public, threadID int

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

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "empty message id",
			},
		})

		return nil
	}

	id, err := strconv.Atoi(values.Get("id"))
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

	text := req.PostFormValue("text")

	if req.PostFormValue("thread") != "" {
		threadID, err = strconv.Atoi(req.PostFormValue("thread"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid thread",
				},
			})

			return
		}
	} else {
		threadID = -1
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

	return h.updateMessage(req.Context(), w, int64(id), user, text, int64(threadID), public)
}

func (h *Handler) updateMessage(ctx context.Context, w http.ResponseWriter, messageID int64, user *users.User,
	text string, threadID int64, public int,
) (err error) {
	var private int32

	if public == 1 {
		private = 0
	} else if public == 0 {
		private = 1
	} else {
		private = -1
	}

	if user.ID == users.PublicUserID && threadID != -1 {
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

	msg, err := h.controller.ReadMessage(ctx, messageID, []int64{user.ID})
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

	if text == "" && threadID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "text or thread_id or both must be provided",
			},
		})

		return nil
	}

	if text == "" {
		text = msg.Text
	}

	err = h.controller.UpdateMessage(ctx, messageID, text, nil, threadID, user.ID, private)
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