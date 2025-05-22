package http

import (
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Logout(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		log.Error("bad cookie")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "bad cookie",
		})

		return err
	}

	token := cookie.Value

	deleteToken(w, h.config.CookieDomain)

	err = h.controller.DeleteToken(req.Context(), log, &model.DeleteTokenParams{
		Token: token,
	})
	if err != nil {
		log.Errorw("failed to delete token, continue", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to delete token",
		})

		return err
	}

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "logged out",
	})
	return nil
}
