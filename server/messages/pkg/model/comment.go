package model

type (
	Comment struct {
		ID           int64          `json:"id"`
		MessageID    int64          `json:"message_id"`
		UserID       int64          `json:"user_id"`
		Text         string         `json:"text"`
		CreatedAt    string         `json:"created_at"`
		UpdatedAt    string         `json:"updated_at"`
		Metadata     []byte         `json:"-"` // open if neccessary
	}

	CommentsList struct {
		Comments            []*Comment
		IsLastPage          bool
		IsFirstPage         bool
		Total               int32
		Offset              int32
		Count               int32
	}

	SendCommentRequest struct {
		MessageID    int64          `json:"message_id"`
		Text         string         `json:"text"`
		Metadata     []byte         `json:"metadata"`
	}

	SendCommentResponse struct {
		ID           int64          `json:"id"`
		Description  string         `json:"description"`
	}

	UpdateCommentRequest struct {
		ID           int64          `json:"id"`
		Text         *string        `json:"text,omitempty"`
	}

	UpdateCommentResponse struct {
		Description  string         `json:"description"`
	}

	DeleteCommentRequest struct {
		ID           int64          `json:"id"`
	}

	DeleteCommentResponse struct {
		Description  string         `json:"description"`
	}

	ReadCommentRequest struct {
		ID           int64          `json:"id"`
	}

	ReadCommentResponse struct {
		Comment      *Comment        `json:"comment"`
	}

	ListCommentsRequest struct {
		MessageID    *int64          `json:"message_id,omitempty"`
		UserID       *int64          `json:"user_id,omitempty"`
		Limit        int32           `json:"limit"`
		Offset       int32           `json:"offset"`
		Asc          bool            `json:"asc"`
	}

	ListCommentsResponse struct {
		Comments     []*Comment   `json:"comments"`
		IsLastPage   bool         `json:"is_last_page"`
		IsFirstPage  bool         `json:"is_first_page"`
		Total        int32        `json:"total"`
		Count        int32        `json:"count"`
		Offset       int32        `json:"offset"`
	}
)