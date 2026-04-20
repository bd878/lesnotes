package http

import (
	"net/http"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) UpdateThread(w http.ResponseWriter, req *http.Request) (err error) {
	var id int64

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
			Error: &server.ErrorCode{
				Code:    threads.CodeThreadPublic,
				Explain: "can update thread of a user only",
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
				Explain: "empty thread id",
			},
		})

		return
	}

	description := req.PostFormValue("description")
	title := req.PostFormValue("title")
	name := req.PostFormValue("name")

	err = h.controller.UpdateThread(req.Context(), id, user.ID, &name, &description, &title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to update thread",
			},
		})

		return err
	}

	response, err := json.Marshal(threads.UpdateThreadResponse{
		ID:          id,
		Description: "updated",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:    "ok",
		Response:  json.RawMessage(response),
	})

	return
}
