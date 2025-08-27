package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) LoginJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:     server.CodeWrongFormat,
				Explain: "failed to parse request",
			},
		})

		return err
	}

	var request users.LoginRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse login request",
			},
		})

		return
	}

	if request.Login == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:     users.CodeNoLogin,
				Explain: "no login",
			},
		})

		return errors.New("cannot get user login from request")
	}

	if request.Password == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    users.CodeNoPassword,
				Explain: "no password",
			},
		})

		return errors.New("cannot get password from request")
	}


	session, err := h.controller.LoginUser(req.Context(), request.Login, request.Password)
	switch err {
	case controller.ErrUserNotFound:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "no user,password pair",
			},
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error:   &server.ErrorCode{
				Code: users.CodeBadPassword,
				Explain: "wrong password",
			},
		})

		return err

	case nil:
		response, err := json.Marshal(users.LoginResponse{
			Token:          session.Token,
			ExpiresUTCNano: session.ExpiresUTCNano,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:       "ok",
			Response:     json.RawMessage(response),
		})

		return nil

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "cannot get user",
			},
		})

		return err
	}

	return nil
}