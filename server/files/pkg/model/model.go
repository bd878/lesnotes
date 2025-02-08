package model

type File struct {
  ID             int32         `json:"id"`
  Name           string        `json:"name,omitempty"`
  CreateUTCNano  int64         `json:"create_utc_nano,omitempty"`
  Error          string        `json:"error,omitempty"`
}

type ReadFileParams struct {
  ID             int32         `json:"id"`
}