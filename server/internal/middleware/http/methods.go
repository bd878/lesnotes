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
	return handler(func(log *logger.Logger, w http.ResponseWriter, req *http.Request) (err error) {
		log.Infof("--> %s\n", req.URL.String())
		err = next.Handle(log, w, req)
		log.Infof("<-- %s\n", req.URL.String())
		if err != nil {
			log.Errorln(err.Error())
		}
		return
	})
}

type UserContextKey struct {}

type userGateway interface {
	Auth(ctx context.Context, log *logger.Logger, token string) (*usermodel.User, error)
	GetPublicUser(ctx context.Context, log *logger.Logger, id int32) (*usermodel.User, error)
}

type authBuilder struct {
	gateway userGateway
	publicUserID int32
	next Handler
}

func AuthBuilder(gateway userGateway, publicUserID int32) MiddlewareFunc {
	return func(next Handler) Handler {
		return &authBuilder{gateway: gateway, publicUserID: publicUserID, next: next}
	}
}

func (b *authBuilder) Handle(log *logger.Logger, w http.ResponseWriter, req *http.Request) (err error) {
	var (
		user *usermodel.User
		cookie *http.Cookie
	)

	cookie, err = req.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		user, err = b.restorePublicUser(log, w, req)
	case nil:
		user, err = b.restoreAuthorizedUser(log, w, req, cookie)
	default:
		log.Errorw("bad cookie", "cookie", cookie, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "bad cookie",
		})

		return
	}

	if err != nil {
		log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "token not found",
		})
		return
	}

	req = req.WithContext(context.WithValue(context.Background(), UserContextKey{}, user))

	return b.next.Handle(log, w, req)
}

func (b *authBuilder) restorePublicUser(log *logger.Logger, w http.ResponseWriter, req *http.Request) (*usermodel.User, error) {
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

func (b *authBuilder) restoreAuthorizedUser(log *logger.Logger, w http.ResponseWriter, req *http.Request, cookie *http.Cookie) (*usermodel.User, error) {
	if cookie == nil {
		return nil, ErrEmptyToken
	}

	user, err := b.gateway.Auth(req.Context(), log, cookie.Value)
	if err != nil {
		log.Errorw("middleware failed to authorize user, gateway error", "cookie", cookie.Value)
		return nil, ErrNoUser
	}

	return user, nil
}