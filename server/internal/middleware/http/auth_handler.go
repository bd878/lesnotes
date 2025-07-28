package middleware

import (
	"errors"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/logger"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	sessionsmodel "github.com/bd878/gallery/server/sessions/pkg/model"
)

var (
	ErrEmptyToken  = errors.New("empty token")
	ErrNoUser      = errors.New("no user")
	ErrNoSession   = errors.New("no session")
	ErrNoPublicID  = errors.New("no public user id")
)

type UserContextKey struct {}

type UsersGateway interface {
	GetUser(ctx context.Context, userID int32) (*usersmodel.User, error)
}

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (*sessionsmodel.Session, error)
}

type authBuilder struct {
	log           *logger.Logger
	users         UsersGateway
	sessions      SessionsGateway
	publicUserID  int32
	next          Handler
}

func AuthBuilder(log *logger.Logger, users UsersGateway, sessions SessionsGateway, publicUserID int32) MiddlewareFunc {
	return func(next Handler) Handler {
		return &authBuilder{log: log, users: users, sessions: sessions, publicUserID: publicUserID, next: next}
	}
}

func (b *authBuilder) Handle(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		user   *usersmodel.User
		cookie *http.Cookie
	)

	cookie, err = req.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		user, err = b.restorePublicUser(req)
	case nil:
		user, err = b.restoreAuthorizedUser(req, cookie)
	default:
		b.log.Errorw("bad cookie", "cookie", cookie, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "bad cookie",
		})

		return
	}

	if err != nil {
		b.log.Errorw("auth middleware failed to restore user, error occured, exit", "error", err)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "error",
			Description: "token not found",
		})
		return
	}

	req = req.WithContext(context.WithValue(context.Background(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *authBuilder) restorePublicUser(req *http.Request) (user *usersmodel.User, err error) {
	if b.publicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err = b.users.GetUser(req.Context(), b.publicUserID)

	return user, nil
}

func (b *authBuilder) restoreAuthorizedUser(req *http.Request, cookie *http.Cookie) (user *usersmodel.User, err error) {
	if cookie == nil {
		return nil, ErrEmptyToken
	}

	session, err := b.sessions.GetSession(req.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}

	user, err = b.users.GetUser(req.Context(), session.UserID)

	return
}