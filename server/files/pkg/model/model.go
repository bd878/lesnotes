package model

type (
	File struct {
		ID             int64         `json:"id"`
		OID            int32         `json:"-"`
		UserID         int64         `json:"user_id"`
		Name           string        `json:"name"`
		CreateUTCNano  int64         `json:"create_utc_nano"`
		Error          string        `json:"error"`
		Size           int64         `json:"size"`
		Mime           string        `json:"mime"`
		Private        bool          `json:"private"`
		Description    string        `json:"description"`
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

	ReadFileMetaRequest struct {
		ID             int64        `json:"id"`
	}

	PrivateFileRequest struct {
		ID             int64        `json:"id"`
	}

	PublishFileRequest struct {
		ID             int64        `json:"id"`
	}

	DeleteFileRequest struct {
		ID             int64        `json:"id"`
	}

	UploadResponse struct {
		ID             int64         `json:"id"`
		Name           string        `json:"name"`
		Description    string        `json:"description"`
	}

	ListFilesResponse struct {
		Files               []*File            `json:"files"`
		IsLastPage          bool               `json:"is_last_page"`
		IsFirstPage         bool               `json:"is_first_page"`
		Count               int32              `json:"count"`
		Total               int32              `json:"total"`
		Offset              int32              `json:"offset"`
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

	DeleteFileResponse struct {
		ID                  int64              `json:"id"`
		Description         string             `json:"description"`
	}

	ReadFileMetaResponse struct {
		File                *File              `json:"file"`
	}
)