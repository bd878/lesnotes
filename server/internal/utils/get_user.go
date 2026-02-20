package utils

import (
	"net/http"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func GetUser(w http.ResponseWriter, req *http.Request) (*users.User, bool) {
	user, ok := req.Context().Value(middleware.UserContextKey{}).(*users.User)
	if !ok {
		return nil, false
	}
	return user, true
}