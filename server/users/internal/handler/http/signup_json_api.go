package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) SignupJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	var data []byte

	data, err = io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
			},
		})

		return
	}

	var body users.SignupRequest
	if err = json.Unmarshal(data, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{	
			Status: "error",
			Error:   &server.ErrorCode {
				Code:     server.CodeWrongFormat,
				Explain: "failed to parse signup body",
			},
		})

		return
	}

	if body.Login == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code: users.CodeNoLogin,
				Explain: "no login",
			},
		})

		return errors.New("cannot get user login from request")
	}

	if body.Password == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:    users.CodeNoPassword,
				Explain: "no password",
			},
		})

		return errors.New("cannot get password from request")
	}

	eightOrMore, twoLetters, oneNumber, oneSpecial := verifyPassword(body.Password)
	if !eightOrMore {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:  users.CodePasswordTooShort,
				Explain: "password is less than 8 symbols",
			},
		})
		return errors.New("password must have > 8 symbols")
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
		return errors.New("upper and lower letters required")
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
		return errors.New("password must have at least one number")
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
		return errors.New("password must have at least one special symbol")
	}

	id := utils.RandomID()

	user, err := h.controller.CreateUser(req.Context(), int64(id), body.Login, body.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:    users.CodeRegisterFailed,
				Explain: "cannot signup user",
			},
		})

		return err
	}

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