package model

type User struct {
  ID               int32            `json:"id"`
  Name             string           `json:"name"`
  Password         string           `json:"password"`
  Token            string           `json:"token"`
  ExpiresUTCNano   int64            `json:"expires_utc_nano"`
}

type ServerResponse struct {
  Status           string           `json:"status"`
  Description      string           `json:"description"`
}

type ServerAuthorizeResponse struct {
  ServerResponse
  Expired          bool             `json:"expired"`
  User             User             `json:"user"`
}
