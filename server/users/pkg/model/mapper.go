package model

import "github.com/bd878/gallery/server/api/users"

func UserToProto(u *User) *users.User {
	return &users.User{
		Id:               u.ID,
		Login:            u.Login,
		HashedPassword:   u.HashedPassword,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		Metadata:         u.Metadata,
		IsPremium:        u.IsPremium,
	}
}

func UserFromProto(u *users.User) *User {
	return &User{
		ID:               u.Id,
		Login:            u.Login,
		HashedPassword:   u.HashedPassword,
		Metadata:         u.Metadata,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		IsPremium:        u.IsPremium,
	}
}