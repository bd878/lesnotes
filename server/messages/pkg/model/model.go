package model

import (
	files "github.com/bd878/gallery/server/files/pkg/model"
)

type (
	SaveMessageResult struct {
		ID                int64
		CreateUTCNano     int64
		UpdateUTCNano     int64
		Private           bool
	}

	ReadOneMessageParams struct {
		ID int64
		UserIDs []int64
	}

	UpdateMessageParams struct {
		ID                int64
		ThreadID          int64
		UserID            int64
		FileIDs           []int64
		Text              string
		UpdateUTCNano     int64
		Private           int32
	}

	UpdateMessageResult struct {
		ID                int64
		UpdateUTCNano     int64
		Private           bool
	}

	DeleteMessageParams struct {
		ID                int64
		UserID            int64
		FileID            int64
	}

	DeleteMessagesParams struct {
		IDs               []int64
		UserID            int64
	}

	PublishMessagesParams struct {
		IDs               []int64
		UserID            int64
		UpdateUTCNano     int64
	}

	DeleteAllUserMessagesParams struct {
		UserID            int64
	}

	PrivateMessagesParams struct {
		IDs               []int64
		UserID            int64
		UpdateUTCNano     int64
	}

	PublishMessagesResult struct {
		UpdateUTCNano     int64
	}

	PrivateMessagesResult struct {
		UpdateUTCNano     int64
	}

	DeleteMessageResult struct {
	}

	DeleteMessagesResult struct {
		IDs               []*DeleteMessageStatus
	}

	ReadThreadMessagesParams struct {
		UserID            int64
		ThreadID          int64
		Limit             int32
		Offset            int32
		Ascending         bool
		Private           int32
	}

	ReadThreadMessagesResult struct {
		Messages          []*Message
		IsLastPage        bool
	}

	ReadMessagesParams struct {
		UserID            int64
		Limit             int32
		Offset            int32
		Ascending         bool
		Private           int32
	}

	ReadMessagesResult struct {
		Messages          []*Message
		IsLastPage        bool
	}

	ReadBatchFilesParams struct {
		UserID            int64
		IDs               []int64
	}

	ReadBatchFilesResult struct {
		Files             map[int64](*files.File)
	}

	SaveFileParams struct {
		Name              string
		UserID            int64
	}

	SaveFileResult struct {
		ID                int64
		CreateUTCNano     int64
	}

	DeleteMessageStatus struct {
		ID                int64           `json:"id"`
		OK                bool            `json:"ok"`
		Explain           string          `json:"explain"`
	}
)
