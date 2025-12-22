package model

type (
	Thread struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
		ParentID    int64   `json:"parent_id"`
		NextID      int64   `json:"next_id"`
		PrevID      int64   `json:"prev_id"`
		Name        string  `json:"name"`
		Count       int32   `json:"count"`
		Private     bool    `json:"private"`
		Description string  `json:"description"`
	}

	ListThreadsRequest struct {
		UserID      int64   `json:"user_id"`
		ParentID    int64   `json:"parent"`
		Limit       int32   `json:"limit"`
		Offset      int32   `json:"offset"`
		Asc         bool    `json:"asc"`
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
		Description string  `json:"description"`
	}

	ReadThreadRequest struct {
		ID          int64   `json:"id"`
		Name        string  `json:"name"`
		UserID      int64   `json:"user_id"`
	}

	UpdateThreadRequest struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
	}

	DeleteThreadRequest struct {
		ID          int64   `json:"id"`
		ParentID    int64   `json:"parent"`
	}

	ResolveThreadRequest struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
	}

	PublishThreadRequest struct {
		ID          int64   `json:"id"`
	}

	PrivateThreadRequest struct {
		ID          int64   `json:"id"`
	}

	ReadThreadResponse struct {
		// TODO: count, offset, total, is_last_page...
		Threads     []*Thread    `json:"threads"`
		Description string       `json:"description"`
	}

	ListThreadsResponse struct {
		Count       int32     `json:"count"`
		IsLastPage  bool      `json:"is_last_page"`
		List        []*Thread `json:"list"`
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

	UpdateThreadResponse struct {
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