package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) PublishMessages(w http.ResponseWriter, req *http.Request) (err error) {
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

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeMessagePublic,
				Explain: "can publish messages of a user only",
			},
		})

		return
	}

	values := req.URL.Query()

	if values.Get("ids") != "" {
		var ids []int64
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: "wrong \"ids\" query field format",
				},
			})

			return err
		}

		return h.publishMessages(w, req, user, ids)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "ids query param required",
			},
		})

		return
	}
}

func (h *Handler) publishMessages(w http.ResponseWriter, req *http.Request, user *users.User, ids []int64) (err error) {
	_, err = h.controller.PublishMessages(req.Context(), &messages.PublishMessagesParams{
		IDs:    ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodePublishFailed,
				Explain: "failed to publish messages",
			},
		})

		return
	}

	response, err := json.Marshal(messages.PublishResponse{
		IDs:           ids,
		Description:   "published",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return
}