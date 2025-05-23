package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	messagesmodel "github.com/bd878/gallery/server/messages/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DeleteJsonAPI(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
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

	if jsonRequest.Name == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no name",
		})

		return errors.New("cannot get name from request")
	}

	if jsonRequest.Password == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{Token: jsonRequest.Token})
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.DeleteUserJsonServerResponse{
			ServerResponse: servermodel.ServerResponse{
				Status: "error",
				Description: "token expired",
			},
			Expired: true,
		})

		return err
	}

	err = h.controller.DeleteUser(req.Context(), log, &model.DeleteUserParams{
		ID: user.ID,
		Token: jsonRequest.Token,
		Name: user.Name,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot delete user",
		})

		return err
	}

	err = h.messagesGateway.DeleteAllUserMessages(req.Context(), log, &messagesmodel.DeleteAllUserMessagesParams{
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot delete all user messages",
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