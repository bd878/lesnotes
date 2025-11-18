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

	ReorderThreadRequest struct {
		ID          int64   `json:"id"`
		ParentID    int64   `json:"parent"`
		NextID      int64   `json:"next"`
		PrevID      int64   `json:"prev"`
	}

	CreateThreadRequest struct {
		ID          int64   `json:"id"`
		ParentID    int64   `json:"parent"`
		NextID      int64   `json:"next"`
		PrevID      int64   `json:"prev"`
		Name        string  `json:"name"`
		Private     bool    `json:"private"`
	}

	DeleteThreadRequest struct {
		ID          int64   `json:"id"`
		ParentID    int64   `json:"parent"`
	}

	ResolveThreadRequest struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
	}

	ResolveThreadResponse struct {
		Path        []int64 `json:"path"`
	}

	ReorderThreadResponse struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
	}

	CreateThreadResponse struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
	}

	DeleteThreadResponse struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
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