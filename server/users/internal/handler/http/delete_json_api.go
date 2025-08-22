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

func (h *Handler) DeleteJsonAPI(w http.ResponseWriter, req *http.Request) error {
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

	var jsonRequest model.DeleteUserJsonRequest
	if err := json.Unmarshal(data, &jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse json request",
		})

		return err
	}

	if jsonRequest.Token == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no token",
		})

		return errors.New("cannot get token from request")
	}

	if jsonRequest.Login == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no login",
		})

		return errors.New("cannot get login from request")
	}

	if jsonRequest.Password == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return errors.New("user required")
	}

	err = h.controller.DeleteUser(req.Context(), int64(user.ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot delete user",
		})

		return err
	}

	json.NewEncoder(w).Encode(&model.DeleteUserJsonServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "deleted",
		},
	})

	return nil
}