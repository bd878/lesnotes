package model

import (
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type SaveMessageParams struct {
  Message          *Message
}

type UpdateMessageParams struct {
  Message          *Message
}

type ReadUserMessagesParams struct {
  UserID            usermodel.UserId
  Limit             int32
  Offset            int32
  Ascending         bool
}

type AuthParams struct {
  Token             string
}

type PutParams struct {
  Message          *Message
}

type GetParams struct {
  UserID            usermodel.UserId
  Limit             int32
  Offset            int32
  Ascending         bool
}