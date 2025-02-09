package model

type File struct {
  ID             int32         `json:"id"`
  UserID         int32         `json:"user_id"`
  Name           string        `json:"name,omitempty"`
  CreateUTCNano  int64         `json:"create_utc_nano,omitempty"`
  Error          string        `json:"error,omitempty"`
}

type ReadFileParams struct {
  ID             int32         `json:"id"`
  UserID         int32         `json:"user_id"`
}

type ReadFileStreamParams struct {
  FileID         int32         `json:"file_id"`
  UserID         int32         `json:"user_id"`
}

type ServerResponse struct {
  Status         string        `json:"status"`
  Description    string        `json:"description"`
}