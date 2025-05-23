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

type DeleteUserParams struct {
	ID int32
	Name string
	Token string
}

type GetUserParams struct {
	ID  int32
	Name string
	Token string
}