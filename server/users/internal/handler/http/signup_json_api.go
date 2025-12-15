package http

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/bd878/gallery/server/i18n"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/users/internal/controller"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) SignupJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	var data []byte

	lang, ok := req.Context().Value(middleware.LangContextKey{}).(i18n.LangCode)
	if !ok {
		lang = i18n.LangEn
	}

	data, err = io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
				Human:   lang.Error(server.CodeWrongFormat),
			},
		})

		return
	}

	var request users.SignupRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{	
			Status: "error",
			Error:   &server.ErrorCode {
				Code:     server.CodeWrongFormat,
				Explain: "failed to parse signup request",
				Human:   lang.Error(server.CodeWrongFormat),
			},
		})

		return
	}

	if request.Login == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:    users.CodeNoLogin,
				Explain: "no login",
				Human:   lang.Error(users.CodeNoLogin),
			},
		})

		return
	}

	if request.Password == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:    users.CodeNoPassword,
				Explain: "no password",
				Human:   lang.Error(users.CodeNoPassword),
			},
		})

		return
	}

	fiveOrMore, _, _, _ := verifyPassword(request.Password)
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

	id := utils.RandomID()

	var user *users.User
	user, err = h.controller.CreateUser(req.Context(), int64(id), request.Login, request.Password)
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
				Error:   &server.ErrorCode{
					Code:    users.CodeSignupFailed,
					Explain: "cannot signup user",
					Human:   lang.Error(users.CodeSignupFailed),
				},
			})
		}

		return
	}

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