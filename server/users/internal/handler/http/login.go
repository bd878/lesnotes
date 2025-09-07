package http

import (
	"time"
	"net/http"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) (err error) {
	var login, password string

	err = req.ParseMultipartForm(1024 /* 1 KB */)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  server.CodeNoForm,
				Explain: "failed to parse form",
			},
		})

		return
	}

	login, password = req.PostFormValue("login"), req.PostFormValue("password")

	if login == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodeNoLogin,
				Explain: "login required",
			},
		})

		return errors.New("login required")
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodeNoPassword,
				Explain: "password required",
			},
		})

		return errors.New("password required")
	}

	session, err := h.controller.LoginUser(req.Context(), login, password)
	switch err {
	case controller.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:     users.CodeAuthFailed,
				Explain: "no login,password pair",
			},
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongPassword,
				Explain: "wrong password",
			},
		})

		return err

	case nil:
		// attach session to response
		http.SetCookie(w, &http.Cookie{
			Name:          "token",
			Value:          session.Token,
			Domain:         h.config.CookieDomain,
			Expires:        time.Unix(0, session.ExpiresUTCNano),
			Path:          "/",
			HttpOnly:       true,
		})

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "cannot get user",
			},
		})

		return err
	}

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
}
