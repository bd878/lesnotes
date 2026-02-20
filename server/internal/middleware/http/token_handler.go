package middleware

import (
	"io"
	"errors"
	"bytes"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/internal/logger"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

type RequestContextKey struct {}

type tokenAuthBuilder struct {
	log             *logger.Logger
	users           UsersGateway
	sessions        SessionsGateway
	publicUserID    int64
	next            Handler
}

func TokenAuthBuilder(log *logger.Logger, users UsersGateway, sessions SessionsGateway, publicUserID int64) MiddlewareFunc {
	return func(next Handler) Handler {
		return &tokenAuthBuilder{log: log, users: users, sessions: sessions, publicUserID: publicUserID, next: next}
	}
}

func (b *tokenAuthBuilder) Handle(w http.ResponseWriter, req *http.Request) (err error) {
	if bytes.Contains([]byte(req.Header.Get("content-type")), []byte("multipart/form-data")) {
		return b.handleMultipartFormData(w, req)
	} else {
		return b.handleJson(w, req)
	}
}

func (b *tokenAuthBuilder) handleMultipartFormData(w http.ResponseWriter, req *http.Request) (err error) {
	var user *users.User

	if req.URL.Query().Has("token") {
		user, err = b.restoreAuthorizedUser(w, req, req.URL.Query().Get("token"))
	} else {
		err = errors.New("multipart/form-data has no token param")
	}

	if err != nil {
		logger.Errorln(err)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoToken,
				Explain: "token not found",
			},
		})
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *tokenAuthBuilder) handleJson(w http.ResponseWriter, req *http.Request) (err error) {
	var user *users.User

	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read body",
			},
		})

		return err
	}

	var request server.ServerRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
			},
		})

		return err
	}

	if request.Token != "" {
		user, err = b.restoreAuthorizedUser(w, req, request.Token)
	} else {
		user, err = b.restorePublicUser(w, req)
	}

	if err != nil {
		logger.Errorln(err)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoToken,
				Explain: "token not found",
			},
		})
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), RequestContextKey{}, request.Request))
	req = req.WithContext(context.WithValue(req.Context(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *tokenAuthBuilder) restorePublicUser(w http.ResponseWriter, req *http.Request) (*users.User, error) {
	if b.publicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err := b.users.GetUser(req.Context(), b.publicUserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *tokenAuthBuilder) restoreAuthorizedUser(w http.ResponseWriter, req *http.Request, token string) (user *users.User, err error) {
	session, err := b.sessions.GetSession(req.Context(), token)
	if err != nil {
		return nil, ErrNoSession
	}

	user, err = b.users.GetUser(req.Context(), int64(session.UserID))

	return
}