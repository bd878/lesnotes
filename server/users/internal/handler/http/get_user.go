package http

import (
	"net/http"
	"strconv"
	"errors"
	"encoding/json"

	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) GetUser(w http.ResponseWriter, req *http.Request) error {
	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "empty user id",
			},
		})

		return errors.New("empty user id")
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "invalid id",
			},
		})

		return err
	}

	user, err := h.controller.GetUser(req.Context(), int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "cannot find user",
			},
		})

		return err
	}

	response, err := json.Marshal(users.GetUserResponse{
		User: user,
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
