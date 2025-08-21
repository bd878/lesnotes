package model

import (
	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type SaveMessageResult struct {
	ID                int64
	CreateUTCNano     int64
	UpdateUTCNano     int64
	Private           bool
}

type ReadOneMessageParams struct {
	ID int64
	UserIDs []int64
}

type UpdateMessageParams struct {
	ID                int64
	ThreadID          int64
	UserID            int64
	FileID            int64
	Text              string
	UpdateUTCNano     int64
	Private           int32
}

type UpdateMessageResult struct {
	ID                int64
	UpdateUTCNano     int64
	Private           bool
}

type DeleteMessageParams struct {
	ID                int64
	UserID            int64
	FileID            int64
}

type DeleteMessagesParams struct {
	IDs               []int64
	UserID            int64
}

type PublishMessagesParams struct {
	IDs               []int64
	UserID            int64
	UpdateUTCNano     int64
}

type DeleteAllUserMessagesParams struct {
	UserID            int64
}

type PrivateMessagesParams struct {
	IDs               []int64
	UserID            int64
	UpdateUTCNano     int64
}

type PublishMessagesResult struct {
	UpdateUTCNano     int64
}

type PrivateMessagesResult struct {
	UpdateUTCNano     int64
}

type DeleteMessageResult struct {
}

type DeleteMessagesResult struct {
	IDs               []*DeleteMessageStatus
}

type ReadThreadMessagesParams struct {
	UserID            int64
	ThreadID          int64
	Limit             int32
	Offset            int32
	Ascending         bool
	Private           int32
}

type ReadThreadMessagesResult struct {
	Messages          []*Message
	IsLastPage        bool
}

type ReadMessagesParams struct {
	UserID            int64
	Limit             int32
	Offset            int32
	Ascending         bool
	Private           int32
}

type ReadMessagesResult struct {
	Messages          []*Message
	IsLastPage        bool
}

type ReadBatchFilesParams struct {
	UserID            int64
	IDs               []int64
}

type ReadBatchFilesResult struct {
	Files             map[int64](*filesmodel.File)
}

type SaveFileParams struct {
	Name              string
	UserID            int64
}

type SaveFileResult struct {
	ID                int64
	CreateUTCNano     int64
}

type (
	ReadMessageOrMessagesJsonRequest struct {
		Public     *int      `json:"public,omitempty"`
		MessageID  int64     `json:"id"`
		ThreadID   int64     `json:"thread_id"`
		Limit      int       `json:"limit"`
		Offset     int       `json:"offset"`
		Asc        int       `json:"asc"`
	}

	PublishMessageOrMessagesJsonRequest struct {
		MessageID *int64     `json:"id,omitempty"`
		IDs       *[]int64   `json:"ids,omitempty"`
	}

	PrivateMessageOrMessagesJsonRequest struct {
		MessageID *int64     `json:"id,omitempty"`
		IDs       *[]int64   `json:"ids,omitempty"`
	}

	UpdateMessageJsonRequest struct {
		MessageID  *int64    `json:"id,omitempty"`
		ThreadID   *int64    `json:"thread_id,omitempty"`
		FileID     *int64    `json:"file_id,omitempty"`
		Text       *string   `json:"text,omitempty"`
		Public     *int      `json:"public,omitempty"`
	}
)