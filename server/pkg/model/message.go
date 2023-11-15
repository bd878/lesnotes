package model

// This message handler passes to repository
type Message struct {
  Value string `json:"value"`
  File string `json:"file"`
}
