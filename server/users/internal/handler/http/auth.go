package http

import (
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

// TODO: Invalidate stale sessions.
// User logs in in one device, get token,
// then logs in in another device, receives new token.
// Old token invalidates, though not expired...
// Check if stage.lesnotes.space tokens influences on
// lesnotes.space (it has .lesnotes.space domain)
func (h *Handler) Auth(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no token",
		})

		return err
	}

	token := cookie.Value

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{Token: token})
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
			ServerResponse: servermodel.ServerResponse{
				Status: "error",
				Description: "token expired",
			},
			Expired: true,
		})

		return err
	}

	if err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user not found",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
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
