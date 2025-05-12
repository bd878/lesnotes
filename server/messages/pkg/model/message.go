package model

import (
	"github.com/bd878/gallery/server/pkg/model"
	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
)

type Message struct {
	ID                  int32               `json:"id,omitempty"`
	ThreadID            int32               `json:"thread_id,omitempty"`
	CreateUTCNano       int64               `json:"create_utc_nano,omitempty"`
	UpdateUTCNano       int64               `json:"update_utc_nano,omitempty"`
	UserID              int32               `json:"user_id,omitempty"`
	FileID              int32               `json:"file_id,omitempty"`
	File                *filesmodel.File    `json:"file,omitempty"`
	Text                string              `json:"text"`
	Private             bool                `json:"private,omitempty"`
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
	Private             bool          `json:"private"`
}

type DeleteMessageServerResponse struct {
	model.ServerResponse
	ID                  int32         `json:"id"`
}

type PublishMessagesServerResponse struct {
	model.ServerResponse
	IDs                  []int32       `json:"ids"`
	UpdateUTCNano        int64         `json:"update_utc_nano"`
}

type PrivateMessagesServerResponse struct {
	model.ServerResponse
	IDs                  []int32       `json:"ids"`
	UpdateUTCNano        int64         `json:"update_utc_nano"`
}

type DeleteMessageStatus struct {
	ID                int32           `json:"id"`
	OK                bool            `json:"ok"`
	Explain           string          `json:"explain"`
}

type DeleteMessagesServerResponse struct {
	model.ServerResponse
	IDs              []*DeleteMessageStatus `json:"ids"`
}

type MessagesListServerResponse struct {
	model.ServerResponse
	ThreadID            int32         `json:"thread_id"`
	Messages            []*Message    `json:"messages"`
	IsLastPage          bool          `json:"is_last_page"`
}
