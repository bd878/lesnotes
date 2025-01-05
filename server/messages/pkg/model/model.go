package model

import (
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

type SaveMessageParams struct {
  Message *Message
}

type ReadUserMessagesParams struct {
  UserId    usermodel.UserId
  Limit     int32
  Offset    int32
  Ascending bool
}

type AuthParams struct {
  Token     string
}

type PutParams struct {
  Message *Message
}

type GetParams struct {
  UserId    usermodel.UserId
  Limit     int32
  Offset    int32
  Ascending bool
}

type FindByIndexParams struct {
  LogIndex uint64
  LogTerm  uint64
}

type PutBatchParams struct {
  MessagesList []*Message
}

type GetOneParams struct {
  UserId    usermodel.UserId
  MessageId MessageId
}