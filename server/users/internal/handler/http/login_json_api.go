package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) LoginJsonAPI(w http.ResponseWriter, req *http.Request) error {
	var err error

	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to parse request",
		})

		return err
	}

	defer req.Body.Close()

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to parse user",
		})

		return err
	}

	if user.Name == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no name",
		})

		return errors.New("cannot get user name from request")
	}

	if user.Password == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}


	session, err := h.controller.LoginUser(req.Context(), user.Name, user.Password)
	switch err {
	case controller.ErrUserNotFound:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no user,password pair",
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "wrong password",
		})

		return err

	case nil:
		json.NewEncoder(w).Encode(&model.LoginJsonUserServerResponse{
			ServerResponse: servermodel.ServerResponse{
				Status: "ok",
				Description: "authenticated",
			},
			Token:          session.Token,
			ExpiresUTCNano: session.ExpiresUTCNano,
		})
		return nil

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "cannot get user",
		})

		return err
	}

	return nil
}