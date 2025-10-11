package model

type (
	Message struct {
		ID      int64             `json:"id"`
		UserID  int64             `json:"user_id"`
		Text    string            `json:"text"`
		Title   string            `json:"title"`
		Name    string            `json:"name"`
	}

	File struct {
		ID      int64
		UserID  int64
		Name    string
	}

	SearchMessagesRequest struct {
		UserID      int64         `json:"user_id"`
		Query       string        `json:"query"`
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