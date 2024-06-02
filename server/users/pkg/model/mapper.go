package model

import (
  "github.com/bd878/gallery/server/api"
)

func UserToProto(u *User) *api.User {
  return &api.User{
    Id: int32(u.Id),
    Name: u.Name,
    Token: u.Token,
    Expires: u.Expires,
  }
}

func UserFromProto(u *api.User) *User {
  return &User{
    Id: int(u.Id),
    Name: u.Name,
    Token: u.Token,
    Expires: u.Expires,
  }
}