package middleware

import (
	"io"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type RequestContextKey struct {}

type tokenAuthBuilder struct {
	gateway userGateway
	publicUserID int32
	next Handler
}

func TokenAuthBuilder(gateway userGateway, publicUserID int32) MiddlewareFunc {
	return func(next Handler) Handler {
		return &tokenAuthBuilder{gateway: gateway, publicUserID: publicUserID, next: next}
	}
}

func (b *tokenAuthBuilder) Handle(log *logger.Logger, w http.ResponseWriter, req *http.Request) (err error) {
	var user *usermodel.User

	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Errorln("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "failed to read body",
		})

		return err
	}

	var jsonApiRequest model.JSONServerRequest
	if err := json.Unmarshal(data, &jsonApiRequest); err != nil {
		log.Errorln("failed to parse json body request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "failed to parse request",
		})

		return err
	}

	if jsonApiRequest.Token != "" {
		user, err = b.restoreAuthorizedUser(log, w, req, jsonApiRequest.Token)
	} else {
		user, err = b.restorePublicUser(log, w, req)
	}

	if err != nil {
		log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "token not found",
		})
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), UserContextKey{}, user))
	req = req.WithContext(context.WithValue(req.Context(), RequestContextKey{}, jsonApiRequest.Req))

	return b.next.Handle(log, w, req)
}

func (b *tokenAuthBuilder) restorePublicUser(log *logger.Logger, w http.ResponseWriter, req *http.Request) (*usermodel.User, error) {
	if b.publicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err := b.gateway.GetPublicUser(req.Context(), log, b.publicUserID)
	if err != nil {
		log.Errorw("middleware failed to restore public user, gateway error", "id", b.publicUserID)
		return nil, err
	}

	return user, nil
}

func (b *tokenAuthBuilder) restoreAuthorizedUser(log *logger.Logger, w http.ResponseWriter, req *http.Request, token string) (*usermodel.User, error) {
	user, err := b.gateway.Auth(req.Context(), log, token)
	if err != nil {
		log.Errorw("middleware failed to authorize user, gateway error", "cookie", token)
		return nil, ErrNoUser
	}

	return user, nil
}