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

func (h *Handler) ReorderThread(w http.ResponseWriter, req *http.Request) (err error) {
	var id, parent, next, prev int64

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
				Explain: "can publish thread of a user only",
			},
		})

		return
	}

	if req.PostFormValue("id") != "" {
		idStr, err := strconv.Atoi(req.PostFormValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    threads.CodeWrongID,
					Explain: "invalid thread",
				},
			})

			return err
		}

		id = int64(idStr)
	}

	if req.PostFormValue("parent") != "" {
		parentStr, err := strconv.Atoi(req.PostFormValue("parent"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    threads.CodeWrongID,
					Explain: "invalid parent",
				},
			})

			return err
		}

		parent = int64(parentStr)
	} else {
		parent = -1
	}

	if req.PostFormValue("next") != "" {
		nextStr, err := strconv.Atoi(req.PostFormValue("next"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    threads.CodeWrongID,
					Explain: "invalid next",
				},
			})

			return err
		}

		next = int64(nextStr)
	}

	if req.PostFormValue("prev") != "" {
		prevStr, err := strconv.Atoi(req.PostFormValue("prev"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    threads.CodeWrongID,
					Explain: "invalid prev",
				},
			})

			return err
		}

		prev = int64(prevStr)
	}


	err = h.controller.ReorderThread(req.Context(), id, user.ID, parent, next, prev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to reorder thread",
			},
		})

		return err
	}

	response, err := json.Marshal(threads.ReorderThreadResponse{
		ID:          id,
		Description: "reordered",
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