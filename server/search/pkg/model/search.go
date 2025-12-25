package model

type (
	Message struct {
		ID          int64         `json:"id"`
		UserID      int64         `json:"user_id"`
		Text        string        `json:"text"`
		Title       string        `json:"title"`
		Name        string        `json:"name"`
		Private     bool          `json:"private"`
	}

	File struct {
		ID          int64
		UserID      int64
		Name        string
	}

	Thread struct {
		ID          int64         `json:"id"`
		UserID      int64         `json:"user_id"`
		ParentID    int64         `json:"parent_id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Private     bool          `json:"private"`
	}

	SearchMessagesRequest struct {
		UserID      int64         `json:"user_id"`
		Query       string        `json:"query"`
		Public      int           `json:"public"`
		ThreadID    int64         `json:"thread"`
	}

	SearchMessagesResponse struct {
		Messages    []*Message    `json:"list"`
		Count       int32         `json:"count"`
	}

	SearchFilesResponse struct {
		Files       []*File       `json:"list"`
		Count       int32         `json:"count"`
	}
)