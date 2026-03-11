package model

const (
	PublicUserID int64 = 9_999_999
)

type (
	User struct {
		ID               int64            `json:"id"`
		Login            string           `json:"login,omitempty"`
		Token            string           `json:"token,omitempty"`
		HashedPassword   string           `json:"-"`
		ExpiresAt        string           `json:"expires_at,omitempty"` // TODO: remove
		CreatedAt        string           `json:"created_at"`
		UpdatedAt        string           `json:"updated_at"`
		IsPremium        bool             `json:"is_premium"`
		Metadata         []byte           `json:"metadata,omitempty"`
	}

	SignupResponse struct {
		Description      string           `json:"description"`
		ID               int64            `json:"id"`
		Token            string           `json:"token"`
		ExpiresAt        string           `json:"expires_at"`
	}

	SignupRequest struct {
		Login            string           `json:"login"`
		Password         string           `json:"password"`
	}

	LogoutResponse struct {
		Description      string           `json:"description"`
	}

	LoginRequest struct {
		Login            string           `json:"login"`
		Password         string           `json:"password"`
	}

	LoginResponse struct {
		Token            string           `json:"token"`
		ExpiresAt        string           `json:"expires_at"`
	}

	GetMeResponse struct {
		ID               int64            `json:"id"`
		Login            string           `json:"login"`
		IsPremium        bool             `json:"is_premium"`
		CreatedAt        string           `json:"created_at"`
		UpdatedAt        string           `json:"updated_at"`
		Metadata         []byte           `json:"metadata,omitempty"`
	}

	DeleteMeRequest struct {
		Login            string           `json:"login"`
	}

	DeleteMeResponse struct {
		Description      string           `json:"description"`
	}

	AuthResponse struct {
		Expired          bool             `json:"expired"`
	}

	UpdateRequest struct {
		Login            *string          `json:"login,omitempty"`
		Metadata         []byte           `json:"metadata,omitempty"`
	}

	UpdateResponse struct {
		Description      string           `json:"description"`
	}
)