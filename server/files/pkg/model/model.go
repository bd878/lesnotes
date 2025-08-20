package model

import "github.com/bd878/gallery/server/pkg/model"

type File struct {
	ID             int32         `json:"id"`
	OID            int32         `json:"-"`
	UserID         int32         `json:"user_id,omitempty"`
	Name           string        `json:"name,omitempty"`
	CreateUTCNano  int64         `json:"create_utc_nano,omitempty"`
	Error          string        `json:"error,omitempty"`
	Size           int64         `json:"size,omitempty"`
	Mime           string        `json:"mime,omitempty"`
	Private        bool          `json:"private,omitempty"`
}

type ReadFileParams struct {
	ID             int32         `json:"id,omitempty"`
	Name           string        `json:"name,omitempty"`
	UserID         int32         `json:"user_id"`
}

type ReadFileStreamParams struct {
	FileName       string        `json:"name,omitempty"`
	FileID         int32         `json:"file_id,omitempty"`
	UserID         int32         `json:"user_id"`
	Public         bool          `json:"public,omitempty"`
}

type SaveFileParams struct {
	Name              string
	UserID            int32
	Private           bool
}

type SaveFileResult struct {
	ID                int32
	CreateUTCNano     int64
}

type UploadFileServerResponse struct {
	model.ServerResponse
	ID                int32      `json:"id"`
	Name              string     `json:"name,omitempty"`
}