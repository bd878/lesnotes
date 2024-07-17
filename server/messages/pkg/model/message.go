package model

// TODO: type MessageID int
type MessageId int

// This message handler passes to repository
type Message struct {
  Id int `json:"id"`
  CreateTime string `json:"createtime"`
  UserId int `json:"userid"`
  Value string `json:"value"`
  File string `json:"file"`
}

// Response to return to the client
type ServerResponse struct {
  Status string `json:"status"`
  Description string `json:"description"`
}

type NewMessageServerResponse struct {
  ServerResponse
  Message Message `json:"message"`
}
