package model

type (
	Message struct {
		ID                 int64         `json:"id"`
		UserID             int64         `json:"user_id"`
		Text               string        `json:"text"`
		Title              string        `json:"title"`
		Name               string        `json:"name"`
		Private            bool          `json:"private"`
		CreatedAt          string        `json:"created_at"`
		UpdatedAt          string        `json:"updated_at"`
	}

	File struct {
		ID                 int64
		UserID             int64
		Name               string
		CreatedAt          string
		UpdatedAt          string
	}

	Thread struct {
		ID                 int64         `json:"id"`
		UserID             int64         `json:"user_id"`
		ParentID           int64         `json:"parent_id"`
		Name               string        `json:"name"`
		Description        string        `json:"description"`
		Private            bool          `json:"private"`
		CreatedAt          string        `json:"created_at"`
		UpdatedAt          string        `json:"updated_at"`
	}

	Translation struct {
		MessageID          int64         `json:"message"`
		Lang               string        `json:"lang"`
		Title              string        `json:"title"`
		Text               string        `json:"text"`
		CreatedAt          string        `json:"created_at"`
		UpdatedAt          string        `json:"updated_at"`
	}

	SearchMessagesRequest struct {
		Query              string        `json:"query"`
		Public             int           `json:"public"`
		ThreadID           int64         `json:"thread"`
		Name               string        `json:"name"`
	}

	SearchMessagesResponse struct {
		Messages           []*Message    `json:"list"`
		IsLastPage         bool          `json:"is_last_page"`
		IsFirstPage        bool          `json:"is_first_page"`
		Total              int32         `json:"total"`
		Count              int32         `json:"count"`
		Offset             int32         `json:"offset"`
	}

	SearchFilesResponse struct {
		Files              []*File       `json:"list"`
		Count              int32         `json:"count"`
	}
)