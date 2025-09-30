package http

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/bd878/gallery/server/third_party/accept"
	"github.com/bd878/gallery/server/i18n"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) LoginJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	var lang i18n.LangCode

	preferredLang, err := accept.Negotiate(req.Header.Get("Accept-Language"), i18n.AcceptedLangs...)
	if err != nil {
		logger.Errorw("login", "error", err)
		lang = i18n.LangEn
	} else {
		lang = i18n.LangFromString(preferredLang)
	}

	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeWrongFormat)),
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
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeWrongFormat)),
			},
		})

		return
	}

	if request.Login == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    users.CodeNoLogin,
				Explain: "no login",
				Human:   lang.Text(fmt.Sprintf("%d", users.CodeNoLogin)),
			},
		})

		return
	}

	if request.Password == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    users.CodeNoPassword,
				Explain: "no password",
				Human:   lang.Text(fmt.Sprintf("%d", users.CodeNoPassword)),
			},
		})

		return
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
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeNoUser)),
			},
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error:   &server.ErrorCode{
				Code:    server.CodeWrongPassword,
				Explain: "wrong password",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeWrongPassword)),
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
				Code:    server.CodeNoUser,
				Explain: "cannot get user",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeNoUser)),
			},
		})

		return err
	}

	return nil
}