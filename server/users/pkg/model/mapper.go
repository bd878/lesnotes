package model

import (
	"github.com/bd878/gallery/server/api"
)

func UserToProto(u *User) *api.User {
	return &api.User{
		Id:               u.ID,
		Login:            u.Login,
		Theme:            u.Theme,
		Lang:             u.Lang,
		FontSize:         u.FontSize,
		HashedPassword:   u.HashedPassword,
	}
}

func UserFromProto(u *api.User) *User {
	return &User{
		ID:               u.Id,
		Login:            u.Login,
		Theme:            u.Theme,
		Lang:             u.Lang,
		FontSize:         u.FontSize,
		HashedPassword:   u.HashedPassword,
	}
}