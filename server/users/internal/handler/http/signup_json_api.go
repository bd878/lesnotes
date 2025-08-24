package http

import (
	"net/http"
	"io"
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

	var request users.SignupRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{	
			Status: "error",
			Error:   &server.ErrorCode {
				Code:     server.CodeWrongFormat,
				Explain: "failed to parse signup request",
			},
		})

		return
	}

	if request.Login == "" {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code: users.CodeNoLogin,
				Explain: "no login",
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
			},
		})

		return
	}

	eightOrMore, twoLetters, oneNumber, oneSpecial := verifyPassword(request.Password)
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

	id := utils.RandomID()

	user, err := h.controller.CreateUser(req.Context(), int64(id), request.Login, request.Password)
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