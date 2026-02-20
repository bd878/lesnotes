package middleware

import (
	"errors"
	"context"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/internal/logger"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
	sessions "github.com/bd878/gallery/server/sessions/pkg/model"
)

var (
	ErrEmptyToken  = errors.New("empty token")
	ErrNoUser      = errors.New("no user")
	ErrNoSession   = errors.New("no session")
	ErrNoPublicID  = errors.New("no public user id")
)

type UserContextKey struct {}

type UsersGateway interface {
	GetUser(ctx context.Context, userID int64) (*users.User, error)
}

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (*sessions.Session, error)
}

type authBuilder struct {
	log           *logger.Logger
	users         UsersGateway
	sessions      SessionsGateway
	publicUserID  int64
	next          Handler
}

func AuthBuilder(log *logger.Logger, users UsersGateway, sessions SessionsGateway, publicUserID int64) MiddlewareFunc {
	return func(next Handler) Handler {
		return &authBuilder{log: log, users: users, sessions: sessions, publicUserID: publicUserID, next: next}
	}
}

// TODO: add Authorization: bearer <token> auth handler method
// now we have Cookie method and post
func (b *authBuilder) Handle(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		user   *users.User
		cookie *http.Cookie
	)

	cookie, err = req.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		user, err = b.restorePublicUser(req)
	case nil:
		user, err = b.restoreAuthorizedUser(req, cookie)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongToken,
				Explain: "bad cookie",
			},
		})

		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoToken,
				Explain: "token not found",
			},
		})
		return
	}

	req = req.WithContext(context.WithValue(context.Background(), UserContextKey{}, user))

	return b.next.Handle(w, req)
}

func (b *authBuilder) restorePublicUser(req *http.Request) (user *users.User, err error) {
	if b.publicUserID == 0 {
		return nil, ErrNoPublicID
	}

	user, err = b.users.GetUser(req.Context(), b.publicUserID)

	return user, nil
}

func (b *authBuilder) restoreAuthorizedUser(req *http.Request, cookie *http.Cookie) (user *users.User, err error) {
	if cookie == nil {
		return nil, ErrEmptyToken
	}

	session, err := b.sessions.GetSession(req.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}

	user, err = b.users.GetUser(req.Context(), int64(session.UserID))

	return
}