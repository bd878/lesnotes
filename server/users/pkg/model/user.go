package model

// TODO: UserId -> UserID
type UserId int
// TODO: type UserName string

// User model
type User struct {
  Id UserId `json:"id"`
  Name string `json:"name"`
  Password string `json:"password"`
  Token string `json:"token"`
  Expires string `json:"expires"`
}

// Response to return to the client
type ServerResponse struct {
  Status string `json:"status"`
  Description string `json:"description"`
}

type ServerAuthorizeResponse struct {
  ServerResponse
  Expired bool `json:"expired"`
  User User `json:"user"`
}
