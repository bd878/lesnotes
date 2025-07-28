package middleware

import (
	"io"
	"errors"
	"bytes"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/logger"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
)

type RequestContextKey struct {}

type tokenAuthBuilder struct {
	log             *logger.Logger
	users           UsersGateway
	sessions        SessionsGateway
	publicUserID    int32
	next            Handler
}

func TokenAuthBuilder(log *logger.Logger, users UsersGateway, sessions SessionsGateway, publicUserID int32) MiddlewareFunc {
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
	var user *usersmodel.User

	if req.URL.Query().Has("token") {
		user, err = b.restoreAuthorizedUser(w, req, req.URL.Query().Get("token"))
	} else {
		err = errors.New("multipart/form-data has no token param")
	}

	if err != nil {
		b.log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "token not found",
		})
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *tokenAuthBuilder) handleJson(w http.ResponseWriter, req *http.Request) (err error) {
	var user *usersmodel.User

	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		b.log.Errorln("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "failed to read body",
		})

		return err
	}

	var jsonApiRequest model.JSONServerRequest
	if err := json.Unmarshal(data, &jsonApiRequest); err != nil {
		b.log.Errorln("failed to parse json body request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "failed to parse request",
		})

		return err
	}

	if jsonApiRequest.Token != "" {
		user, err = b.restoreAuthorizedUser(w, req, jsonApiRequest.Token)
	} else {
		user, err = b.restorePublicUser(w, req)
	}

	if err != nil {
		b.log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "token not found",
		})
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), RequestContextKey{}, jsonApiRequest.Req))
	req = req.WithContext(context.WithValue(req.Context(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *tokenAuthBuilder) restorePublicUser(w http.ResponseWriter, req *http.Request) (*usersmodel.User, error) {
	if b.publicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err := b.users.GetUser(req.Context(), b.publicUserID)
	if err != nil {
		b.log.Errorw("middleware failed to restore public user, gateway error", "id", b.publicUserID)
		return nil, err
	}

	return user, nil
}

func (b *tokenAuthBuilder) restoreAuthorizedUser(w http.ResponseWriter, req *http.Request, token string) (user *usersmodel.User, err error) {
	session, err := b.sessions.GetSession(req.Context(), token)
	if err != nil {
		return nil, ErrNoSession
	}

	user, err = b.users.GetUser(req.Context(), session.UserID)

	return
}