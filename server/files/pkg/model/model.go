package model

import "github.com/bd878/gallery/server/pkg/model"

type File struct {
	ID             int64         `json:"id"`
	OID            int32         `json:"-"`
	UserID         int64         `json:"user_id,omitempty"`
	Name           string        `json:"name,omitempty"`
	CreateUTCNano  int64         `json:"create_utc_nano,omitempty"`
	Error          string        `json:"error,omitempty"`
	Size           int64         `json:"size,omitempty"`
	Mime           string        `json:"mime,omitempty"`
	Private        bool          `json:"private,omitempty"`
}

type ReadFileParams struct {
	ID             int64         `json:"id,omitempty"`
	Name           string        `json:"name,omitempty"`
	UserID         int64         `json:"user_id"`
}

type ReadFileStreamParams struct {
	FileName       string        `json:"name,omitempty"`
	FileID         int64         `json:"file_id,omitempty"`
	UserID         int64         `json:"user_id"`
	Public         bool          `json:"public,omitempty"`
}

type SaveFileParams struct {
	Name              string
	UserID            int64
	Private           bool
}

type SaveFileResult struct {
	ID                int64
	CreateUTCNano     int64
}

type UploadFileServerResponse struct {
	model.ServerResponse
	ID                int64      `json:"id"`
	Name              string     `json:"name,omitempty"`
}