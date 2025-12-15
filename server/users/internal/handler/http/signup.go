package http

import (
	"time"
	"net/http"
	"encoding/json"
	"unicode"

	"github.com/bd878/gallery/server/i18n"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/users/internal/controller"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func verifyPassword(password string) (fiveOrMore, twoLetters, oneNumber, oneSpecial bool) {
	if len(password) >= 5 {
		fiveOrMore = true
	}

	var oneLower, oneUpper bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			oneNumber = true
		case unicode.IsUpper(c):
			oneUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			oneSpecial = true
		case unicode.IsLetter(c):
			oneLower = true
		default:
		}
	}
	twoLetters = oneUpper && oneLower
	return
}

func verifyLogin(login string) (fiveOrMore bool) {
	if len(login) >= 5 {
		fiveOrMore = true
	}
	return
}

func (h *Handler) Signup(w http.ResponseWriter, req *http.Request) (err error) {
	var login, password string

	lang, ok := req.Context().Value(middleware.LangContextKey{}).(i18n.LangCode)
	if !ok {
		lang = i18n.LangEn
	}

	err = req.ParseMultipartForm(1024 /* 1 KB */)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoForm,
				Explain: "failed to parse form",
				Human:   lang.Error(server.CodeNoForm),
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
				Human:   lang.Error(users.CodeNoLogin),
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
				Human:   lang.Error(users.CodeNoPassword),
			},
		})

		return
	}

	fiveOrMore, _, _, _ := verifyPassword(password)
	if !fiveOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    users.CodePasswordTooShort,
				Explain: "password is less than 5 symbols",
				Human:   lang.Error(users.CodePasswordTooShort),
			},
		})
		return
	}

	fiveOrMore = verifyLogin(login)
	if !fiveOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    users.CodeLoginTooShort,
				Explain: "login is less than 5 symbols",
				Human:   lang.Error(users.CodeLoginTooShort),
			},
		})
		return
	}

	id := utils.RandomID()

	var user *users.User
	user, err = h.controller.CreateUser(req.Context(), int64(id), login, password)
	if err != nil {
		switch err {
		case controller.ErrUserExists:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    users.CodeUserExists,
					Explain: "user exists",
					Human:   lang.Error(users.CodeUserExists),
				},
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    users.CodeSignupFailed,
					Explain: "cannot signup user",
					Human:   lang.Error(users.CodeSignupFailed),
				},
			})
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:          "token",
		Value:          user.Token,
		Domain:         h.config.CookieDomain,
		Expires:        time.Unix(0, user.ExpiresUTCNano),
		Path:          "/",
		HttpOnly:       true,
	})

	response, err := json.Marshal(users.SignupResponse{
		Description:    "user signed up",
		ID:             user.ID,
		Token:          user.Token,
		ExpiresUTCNano: user.ExpiresUTCNano,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:      "ok",
		Response:    json.RawMessage(response),
	})

	return nil
}
