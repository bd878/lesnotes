package model

import (
	"github.com/bd878/gallery/server/api"
)

func UserToProto(u *User) *api.User {
	return &api.User{
		Id:               u.ID,
		Name:             u.Name,
		Token:            u.Token,
		ExpiresUtcNano:   u.ExpiresUTCNano,
	}
}

func UserFromProto(u *api.User) *User {
	return &User{
		ID:               u.Id,
		Name:             u.Name,
		Token:            u.Token,
		ExpiresUTCNano:   u.ExpiresUtcNano,
	}
}