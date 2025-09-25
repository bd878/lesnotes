package model

import (
	"github.com/bd878/gallery/server/api"
)

func UserToProto(u *User) *api.User {
	return &api.User{
		Id:               u.ID,
		Login:            u.Login,
		HashedPassword:   u.HashedPassword,
		Token:            u.Token,
		Theme:            u.Theme,
		ExpiresUtcNano:   u.ExpiresUTCNano,
		Lang:             u.Lang,
		FontSize:         u.FontSize,
	}
}

func UserFromProto(u *api.User) *User {
	return &User{
		ID:               u.Id,
		Login:            u.Login,
		Token:            u.Token,
		Theme:            u.Theme,
		Lang:             u.Lang,
		HashedPassword:   u.HashedPassword,
		ExpiresUTCNano:   u.ExpiresUtcNano,
		FontSize:         u.FontSize,
	}
}