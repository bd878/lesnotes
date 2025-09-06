package model

import (
	files "github.com/bd878/gallery/server/files/pkg/model"
)

type (
	Message struct {
		ID                  int64               `json:"id"`
		ThreadID            int64               `json:"thread"`
		CreateUTCNano       int64               `json:"create_utc_nano,omitempty"`
		UpdateUTCNano       int64               `json:"update_utc_nano,omitempty"`
		UserID              int64               `json:"user_id"`         // TODO: load user, == 0 for public user
		Name                string              `json:"name"`
		FileIDs             []int64             `json:"-"`
		Files               []*files.File       `json:"files"`
		Text                string              `json:"text"`
		Title               string              `json:"title"`
		Private             bool                `json:"private"`
	}

	SendRequest struct {
		Text                string              `json:"text"`
		Title               string              `json:"title"`
		FileIDs             []int64             `json:"file_ids,omitempty"`
		Private             bool                `json:"private"`
		ThreadID            int64               `json:"thread"`
	}

	SendResponse struct {
		Message             *Message            `json:"message"`
	}

	PublishRequest struct {
		IDs                 []int64            `json:"ids"`
	}

	PublishResponse struct {
		IDs                 []int64             `json:"ids"`
		Description         string              `json:"description"`
	}

	PrivateRequest struct {
		IDs                 []int64             `json:"ids"`
	}

	PrivateResponse struct {
		IDs                 []int64             `json:"ids"`
		Description         string              `json:"description"`
	}

	ReadResponse struct {
		ThreadID            *int64              `json:"thread,omitempty"`
		Messages            []*Message          `json:"messages"`
		IsLastPage          *bool               `json:"is_last_page"`
		Description         string              `json:"description"`
	}

	ReadRequest struct {
		UserID              int64               `json:"user"`
		MessageID           int64               `json:"id"`
		ThreadID            int64               `json:"thread"`
		Limit               int                 `json:"limit"`
		Offset              int                 `json:"offset"`
		Asc                 int                 `json:"asc"`
		IDs                 []int64             `json:"ids"`
	}

	ReadPathRequest struct {
		ID                  int64               `json:"id"`
	}

	ReadPathResponse struct {
		Messages            []*Message          `json:"path"`
	}

	DeleteRequest struct {
		IDs                 []int64             `json:"ids"`
	}

	DeleteResponse struct {
		Description         string              `json:"description"`
		IDs                 []int64             `json:"ids"`
	}

	UpdateRequest struct {
		MessageID           int64               `json:"id"`
		ThreadID            *int64              `json:"thread,omitempty"`
		Text                *string             `json:"text,omitempty"`
		Public              *int                `json:"public,omitempty"`
		Title               *string             `json:"title,omitempty"`
	}

	UpdateResponse struct {
		ID                  int64               `json:"id"`
		UpdateUTCNano       int64               `json:"update_utc_nano"`
		Description         string              `json:"description"`
	}
)
