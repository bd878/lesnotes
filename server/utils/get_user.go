package utils

import (
	"net/http"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

func GetUser(w http.ResponseWriter, req *http.Request) (*usermodel.User, bool) {
	user, ok := req.Context().Value(httpmiddleware.UserContextKey{}).(*usermodel.User)
	if !ok {
		return nil, false
	}
	return user, true
}