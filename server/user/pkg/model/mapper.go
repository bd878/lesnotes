package model

import (
  "github.com/bd878/gallery/server/gen"
)

func UserToProto(u *User) *gen.User {
  return &gen.User{
    Id: int32(u.Id),
    Name: u.Name,
    Token: u.Token,
    Expires: u.Expires,
  }
}

func UserFromProto(u *gen.User) *User {
  return &User{
    Id: int(u.Id),
    Name: u.Name,
    Token: u.Token,
    Expires: u.Expires,
  }
}