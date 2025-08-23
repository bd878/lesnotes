package http

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) AuthJsonAPI(w http.ResponseWriter, req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "failed to read json body",
			},
		})

		return err
	}

	var request server.ServerRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "failed to parse json request",
			},
		})

		return err
	}

	user, err := h.controller.AuthUser(req.Context(), request.Token)
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeTokenExpired,
				Explain: "token expired",
			},
		})

		return err
	}

	if err != nil {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "no user",
			},
		})

		return err
	}

	response, err := json.Marshal(users.AuthResponse{
		Expired: false,
		User: user,
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
