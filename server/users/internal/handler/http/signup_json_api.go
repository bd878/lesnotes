package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) SignupJsonAPI(w http.ResponseWriter, req *http.Request) error {
	var err error

	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse request",
		})

		return err
	}

	defer req.Body.Close()

	var body model.SignupUserJsonRequest
	if err := json.Unmarshal(data, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse user",
		})

		return err
	}

	if body.Login == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no name",
		})

		return errors.New("cannot get user name from request")
	}

	if body.Password == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}

	id := utils.RandomID()

	user, err := h.controller.CreateUser(req.Context(), int64(id), body.Login, body.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot check user",
		})

		return err
	}

	json.NewEncoder(w).Encode(&model.SignupJsonUserServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status:      "ok",
			Description: "created",
		},
		Token:          user.Token,
		ExpiresUTCNano: user.ExpiresUTCNano,
		ID:             user.ID,
	})

	return nil
}