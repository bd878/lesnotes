package model

type (
	Thread struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
		ParentID    int64   `json:"parent_id"`
		NextID      int64   `json:"next_id"`
		PrevID      int64   `json:"prev_id"`
		Name        string  `json:"name"`
		Private     bool    `json:"private"`
	}

	ResolveThreadRequest struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
	}

	ResolveThreadResponse struct {
		Path        []int64     `json:"path"`
	}

	PrivateResponse struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
	}

	PublishResponse struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
	}
)