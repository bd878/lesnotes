package http

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) AuthJsonAPI(w http.ResponseWriter, req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to read json body",
		})

		return err
	}

	defer req.Body.Close()

	var jsonReq servermodel.JSONServerRequest
	if err := json.Unmarshal(data, &jsonReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to parse json request",
		})

		return err
	}

	token := jsonReq.Token

	user, err := h.controller.AuthUser(req.Context(), token)
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
			ServerResponse: servermodel.ServerResponse{
				Status:      "error",
				Code:        "expired",
				Description: "token expired",
			},
			Expired: true,
		})

		return err
	}

	if err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Code:        "no_token",
			Description: "user not found",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
		ServerResponse: servermodel.ServerResponse{
			Status:      "ok",
			Description: "token valid",
		},
		Expired: false,
		User: model.User{
			ID:               user.ID,
			Name:             user.Name,
			Token:            user.Token,
			ExpiresUTCNano:   user.ExpiresUTCNano,
		},
	})

	return nil
}
