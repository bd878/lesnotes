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

	List struct {
		Files               []*File
		IsLastPage          bool
		IsFirstPage         bool
		Total               int32
		Count               int32
		Offset              int32
	}

	ListFilesRequest struct {
		UserID         int64         `json:"user"`
		Limit          int32         `json:"limit"`
		Offset         int32         `json:"offset"`
		Asc            int           `json:"asc"`
	}

	PrivateFileRequest struct {
		ID             int64        `json:"id"`
	}

	PublishFileRequest struct {
		ID             int64        `json:"id"`
	}

	UploadResponse struct {
		ID             int64         `json:"id"`
		Name           string        `json:"name,omitempty"`
		Description    string        `json:"description,omitempty"`
	}

	ListFilesResponse struct {
		Files               []*File            `json:"files"`
		IsLastPage          bool               `json:"is_last_page"`
		IsFirstPage         bool               `json:"is_first_page,omitempty"`
		Count               int32              `json:"count,omitempty"`
		Total               int32              `json:"total,omitempty"`
		Offset              int32              `json:"offset,omitempty"`
		Description         string             `json:"description"`
	}

	PrivateFileResponse struct {
		ID                  int64              `json:"id"`
		Description         string             `json:"description"`
	}

	PublishFileResponse struct {
		ID                  int64              `json:"id"`
		Description         string             `json:"description"`
	}
)