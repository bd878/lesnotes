package model

import (
	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type SaveMessageResult struct {
	ID                int32
	CreateUTCNano     int64
	UpdateUTCNano     int64
	Private           bool
}

type UpdateMessageParams struct {
	ID                int32
	ThreadID          int32
	UserID            int32
	FileID            int32
	Text              string
	UpdateUTCNano     int64
	Private           bool
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

type ReadThreadMessagesParams struct {
	UserID            int32
	ThreadID          int32
	Limit             int32
	Offset            int32
	Ascending         bool
}

type ReadThreadMessagesResult struct {
	Messages          []*Message
	IsLastPage        bool
}

type ReadAllMessagesParams struct {
	UserID            int32
	Limit             int32
	Offset            int32
	Ascending         bool
}

type ReadAllMessagesResult struct {
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