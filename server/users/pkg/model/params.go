package model

type AddUserParams struct {
  User      *User
}

type HasUserParams struct {
  User      *User
}

type RefreshTokenParams struct {
  User      *User
}

type DeleteTokenParams struct {
  Token string
}

type GetUserParams struct {
  User      *User
}