package middleware

import (
	"errors"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

var (
	ErrEmptyToken = errors.New("empty token")
	ErrNoUser = errors.New("no user")
	ErrNoPublicID = errors.New("no public user id")
)

func Log(next Handler) Handler {
	return func(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
		next(log, w, req)
	}
}

type UserContextKey struct {}

type userGateway interface {
	Auth(ctx context.Context, log *logger.Logger, token string) (*usermodel.User, error)
	GetPublicUser(ctx context.Context, log *logger.Logger, id int32) (*usermodel.User, error)
}

type AuthBuilder struct {
	Gateway userGateway
	PublicUserID int32
}

func (b *AuthBuilder) Auth(next Handler) Handler {
	return func(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
		var user *usermodel.User
		var err error

		cookie, err := req.Cookie("token")
		switch err {
		case http.ErrNoCookie:
			user, err = b.restorePublicUser(log, w, req)
		case nil:
			user, err = b.restoreAuthorizedUser(log, w, req, cookie)
		default:
			log.Errorw("bad cookie", "cookie", cookie, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
			json.NewEncoder(w).Encode(model.ServerResponse{
				Status: "ok",
				Description: "token not found",
			})
			return
		}

		req = req.WithContext(context.WithValue(context.Background(), UserContextKey{}, user))

		next(log, w, req)
	}
}

func (b *AuthBuilder) restorePublicUser(log *logger.Logger, w http.ResponseWriter, req *http.Request) (*usermodel.User, error) {
	if b.PublicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err := b.Gateway.GetPublicUser(req.Context(), log, b.PublicUserID)
	if err != nil {
		log.Errorw("middleware failed to restore public user, gateway error", "id", b.PublicUserID)
		return nil, err
	}

	return user, nil
}

func (b *AuthBuilder) restoreAuthorizedUser(log *logger.Logger, w http.ResponseWriter, req *http.Request, cookie *http.Cookie) (*usermodel.User, error) {
	if cookie == nil {
		return nil, ErrEmptyToken
	}

	user, err := b.Gateway.Auth(req.Context(), log, cookie.Value)
	if err != nil {
		log.Errorw("middleware failed to authorize user, gateway error", "cookie", cookie.Value)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "token not found",
		})
		return nil, ErrNoUser
	}

	return user, nil
}