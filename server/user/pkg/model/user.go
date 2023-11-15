package model

// User model
type User struct {
  Name string `json:"name"`
  Password string `json:"password"`
  Token string `json:"token"`
}

// Response to return to the client
type ServerResponse struct {
  Status string `json:"status"`
  Description string `json:"description"`
}
