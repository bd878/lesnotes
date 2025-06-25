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

type ReadOneMessageParams struct {
	ID int32
	UserIDs []int32
}

type UpdateMessageParams struct {
	ID                int32
	ThreadID          int32
	UserID            int32
	FileID            int32
	Text              string
	UpdateUTCNano     int64
	Private           int32
}

type UpdateMessageResult struct {
	ID                int32
	UpdateUTCNano     int64
	Private           bool
}

type DeleteMessageParams struct {
	ID                int32
	UserID            int32
	FileID            int32
}

type DeleteMessagesParams struct {
	IDs               []int32
	UserID            int32
}

type PublishMessagesParams struct {
	IDs               []int32
	UserID            int32
	UpdateUTCNano     int64
}

type DeleteAllUserMessagesParams struct {
	UserID            int32
}

type PrivateMessagesParams struct {
	IDs               []int32
	UserID            int32
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
	UserID            int32
	ThreadID          int32
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
	UserID            int32
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
	UserID            int32
	IDs               []int32
}

type ReadBatchFilesResult struct {
	Files             map[int32](*filesmodel.File)
}

type SaveFileParams struct {
	Name              string
	UserID            int32
}

type SaveFileResult struct {
	ID                int32
	CreateUTCNano     int64
}

type (
	ReadMessageOrMessagesJsonRequest struct {
		Public     *int      `json:"public,omitempty"`
		MessageID  int32     `json:"id"`
		ThreadID   int32     `json:"thread_id"`
		Limit      int       `json:"limit"`
		Offset     int       `json:"offset"`
		Asc        int       `json:"asc"`
	}

	PublishMessageOrMessagesJsonRequest struct {
		MessageID *int32     `json:"id,omitempty"`
		IDs       *[]int32   `json:"ids,omitempty"`
	}

	PrivateMessageOrMessagesJsonRequest struct {
		MessageID *int32     `json:"id,omitempty"`
		IDs       *[]int32   `json:"ids,omitempty"`
	}

	UpdateMessageJsonRequest struct {
		MessageID  *int32    `json:"id,omitempty"`
		ThreadID   *int32    `json:"thread_id,omitempty"`
		FileID     *int32    `json:"file_id,omitempty"`
		Text       *string   `json:"text,omitempty"`
		Public     *int      `json:"public,omitempty"`
	}
)