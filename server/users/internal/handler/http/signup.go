package http

import (
	"time"
	"net/http"
	"encoding/json"
	"unicode"

	"github.com/bd878/gallery/server/utils"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func verifyPassword(password string) (eightOrMore, twoLetters, oneNumber, oneSpecial bool) {
	if len(password) >= 8 {
		eightOrMore = true
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

func verifyLogin(login string) (eightOrMore bool) {
	if len(login) >= 8 {
		eightOrMore = true
	}
	return
}

func (h *Handler) Signup( w http.ResponseWriter, req *http.Request) (err error) {
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

		return
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

		return
	}

	eightOrMore, twoLetters, oneNumber, oneSpecial := verifyPassword(password)
	if !eightOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodePasswordTooShort,
				Explain: "password is less than 8 symbols",
			},
		})
		return
	}
	if !twoLetters {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: users.CodePasswordUpperLower,
				Explain: "password must have upper und lower letter",
			},
		})
		return
	}
	if !oneNumber {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodePasswordOneNumber,
				Explain: "password must have at least one number",
			},
		})
		return
	}
	if !oneSpecial {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodePasswordOneSpecial,
				Explain: "password must have at least one special symbol",
			},
		})
		return
	}

	eightOrMore = verifyLogin(login)
	if !eightOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodeLoginTooShort,
				Explain: "login is less than 8 symbols",
			},
		})
		return
	}

	id := utils.RandomID()

	var user *users.User
	user, err = h.controller.CreateUser(req.Context(), int64(id), login, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodeRegisterFailed,
				Explain: "cannot signup user",
			},
		})

		return err
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
		Description:    "user registered",
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
