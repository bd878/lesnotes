package model

import (
	files "github.com/bd878/gallery/server/files/pkg/model"
)

type (
	Message struct {
		ID                  int64               `json:"id,omitempty"`
		ThreadID            int64               `json:"thread_id,omitempty"`
		CreateUTCNano       int64               `json:"create_utc_nano,omitempty"`
		UpdateUTCNano       int64               `json:"update_utc_nano,omitempty"`
		UserID              int64               `json:"user_id,omitempty"`
		FileID              int64               `json:"file_id,omitempty"`
		FileIDs             []int64             `json:"file_ids,omitempty"`
		Name                string              `json:"name,omitempty"`
		File                *files.File         `json:"file,omitempty"`
		Text                string              `json:"text"`
		Private             bool                `json:"private,omitempty"`
	}

	UpdateRequest struct {
		MessageID           int64               `json:"id"`
		ThreadID            *int64              `json:"thread_id,omitempty"`
		FileID              *int64              `json:"file_id,omitempty"`
		Text                *string             `json:"text,omitempty"`
		Public              *int                `json:"public,omitempty"`
	}

	UpdateResponse struct {
		ID                  int64               `json:"id"`
		UpdateUTCNano       int64               `json:"update_utc_nano"`
		Private             bool                `json:"private"`
		Description         string              `json:"description"`
	}

	SaveResponse struct {
		Message             *Message            `json:"message"`
	}

	ReadRequest struct {
		Public              *int                `json:"public,omitempty"`
		MessageID           int64               `json:"id"`
		ThreadID            int64               `json:"thread_id"`
		Limit               int                 `json:"limit"`
		Offset              int                 `json:"offset"`
		Asc                 int                 `json:"asc"`
	}

	ReadResponse struct {
		Message             *Message            `json:"message"`
	}

	ListResponse struct {
		ThreadID            *int64              `json:"thread_id,omitempty"`
		Messages            []*Message          `json:"messages"`
		IsLastPage          bool                `json:"is_last_page"`
	}

	PublishRequest struct {
		MessageID           *int64              `json:"id,omitempty"`
		IDs                 *[]int64            `json:"ids,omitempty"`		
	}

	PublishResponse struct {
		IDs                 []int64             `json:"ids"`
		UpdateUTCNano       int64               `json:"update_utc_nano"`
		Description         string              `json:"description"`
	}

	PrivateRequest struct {
		MessageID           *int64              `json:"id,omitempty"`
		IDs                 *[]int64            `json:"ids,omitempty"`
	}

	PrivateResponse struct {
		IDs                  []int64            `json:"ids"`
		UpdateUTCNano        int64              `json:"update_utc_nano"`
		Description          string             `json:"description"`
	}

	DeleteResponse struct {
		Description          string             `json:"description"`
		ID                   *int64             `json:"id,omitempty"`
		IDs                  *[]int64           `json:"ids,omitempty"`
	}


	DeleteMessageStatus struct {
		ID                int64           `json:"id"`
		OK                bool            `json:"ok"`
		Explain           string          `json:"explain"`
	}
)
