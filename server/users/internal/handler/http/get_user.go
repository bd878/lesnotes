package http

import (
	"net/http"
	"strconv"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) GetUser(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty user id",
		})

		return errors.New("empty user id")
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{ID: int32(id)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot find user",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerUserResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "exists",
		},
		User: model.User{
			ID: user.ID,
			Name: user.Name,
		},
	})

	return nil
}
