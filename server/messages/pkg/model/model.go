package model

import (
  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type SaveMessageParams struct {
  Message          *Message
}

type SaveMessageResult struct {
  ID                int32
  CreateUTCNano     int64
  UpdateUTCNano     int64
}

type UpdateMessageParams struct {
  ID                int32
  UserID            int32
  FileID            int32
  Text              string
  UpdateUTCNano     int64
}

type UpdateMessageResult struct {
  ID                int32
  UpdateUTCNano     int64
}

type DeleteMessageParams struct {
  ID                int32
  UserID            int32
  FileID            int32
}

type DeleteMessageResult struct {
}

type ReadUserMessagesParams struct {
  UserID            int32
  Limit             int32
  Offset            int32
  Ascending         bool
}

type ReadUserMessagesResult struct {
  Messages          []*Message
  IsLastPage        bool
}

type SaveFileParams struct {
  Name              string
  UserID            int32
}

type SaveFileResult struct {
  ID                int32
  CreateUTCNano     int64
}

type ReadBatchFilesParams struct {
  UserID            int32
  IDs               []int32
}

type ReadBatchFilesResult struct {
  Files             map[int32](*filesmodel.File)
}