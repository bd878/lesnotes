package model

type (
	File struct {
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

	UploadResponse struct {
		ID             int64         `json:"id"`
		Name           string        `json:"name,omitempty"`
		Description    string        `json:"description,omitempty"`
	}
)