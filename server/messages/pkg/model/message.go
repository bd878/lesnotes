package model

import (
  "github.com/bd878/gallery/server/pkg/model"
  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type Message struct {
  ID                  int32               `json:"id"`
  CreateUTCNano       int64               `json:"create_utc_nano,omitmempty"`
  UpdateUTCNano       int64               `json:"update_utc_nano,omitmempty"`
  UserID              int32               `json:"user_id,omitmempty"`
  File                *filesmodel.File    `json:"file,omitmempty"`
  Text                string              `json:"text,omitmempty"`
}

type MessagesList struct {
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}

type NewMessageServerResponse struct {
  model.ServerResponse
  Message             Message       `json:"message"`
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
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}
