package model

import (
  "github.com/bd878/gallery/server/pkg/model"
  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type Message struct {
  ID                  int32               `json:"id"`
  ThreadID            int32               `json:"thread_id,omitempty"`
  CreateUTCNano       int64               `json:"create_utc_nano,omitempty"`
  UpdateUTCNano       int64               `json:"update_utc_nano,omitempty"`
  UserID              int32               `json:"user_id,omitempty"`
  File                *filesmodel.File    `json:"file,omitempty"`
  Text                string              `json:"text,omitempty"`
}

type MessagesList struct {
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}

type ReadMessageServerResponse struct {
  model.ServerResponse
  Message             *Message       `json:"message"`
}

type NewMessageServerResponse struct {
  model.ServerResponse
  Message             *Message       `json:"message"`
}

type UpdateMessageServerResponse struct {
  model.ServerResponse
  ID                  int32         `json:"id"`
  UpdateUTCNano       int64         `json:"update_utc_nano"`
}

type DeleteMessageServerResponse struct {
  model.ServerResponse
  ID                  int32         `json:"id"`
}

type MessagesListServerResponse struct {
  model.ServerResponse
  ThreadID            int32         `json:"thread_id,omitempty"`
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}
