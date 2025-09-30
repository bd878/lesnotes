package http

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/third_party/accept"
	"github.com/bd878/gallery/server/i18n"
	"github.com/bd878/gallery/server/logger"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) (err error) {
	var login, password string
	var lang i18n.LangCode

	languages := req.Header.Get("Accept-Language")
	preferredLang, err := accept.Negotiate(languages, i18n.AcceptedLangs...)
	if err != nil {
		logger.Errorw("login", "error", err)
		lang = i18n.LangEn
	} else {
		lang = i18n.LangFromString(preferredLang)
	}

	err = req.ParseMultipartForm(1024 /* 1 KB */)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoForm,
				Explain: "failed to parse form",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeNoForm)),
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
				Code:    users.CodeNoLogin,
				Explain: "login required",
				Human:   lang.Text(fmt.Sprintf("%d", users.CodeNoLogin)),
			},
		})

		return
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    users.CodeNoPassword,
				Explain: "password required",
				Human:   lang.Text(fmt.Sprintf("%d", users.CodeNoPassword)),
			},
		})

		return
	}

	session, err := h.controller.LoginUser(req.Context(), login, password)
	switch err {
	case controller.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:     server.CodeNoUser,
				Explain: "no login,password pair",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeNoUser)),
			},
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongPassword,
				Explain: "wrong password",
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeWrongPassword)),
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
				Human:   lang.Text(fmt.Sprintf("%d", server.CodeNoUser)),
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
