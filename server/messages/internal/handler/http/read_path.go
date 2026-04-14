package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) ReadPath(w http.ResponseWriter, req *http.Request) (err error) {
	var messageID int64

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

	if values.Has("id") {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid message id",
				},
			})

			return err
		}

		messageID = int64(id)
	}

	list, parentID, err := h.controller.ReadPath(req.Context(), user.ID, messageID, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read messages path",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.ReadPathResponse{
		Messages: list,
		ThreadID: parentID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}
