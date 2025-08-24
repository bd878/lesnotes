package model

const (
	PublicUserID int64 = 9_999_999
)

type (
	User struct {
		ID               int64            `json:"id"`
		Login            string           `json:"login,omitempty"`
		HashedPassword   string           `json:"salt,omitempty"`
		Token            string           `json:"token,omitempty"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano,omitempty"` // TODO: ExpiresAt
	}

	SignupResponse struct {
		Description      string           `json:"description"`
		ID               int64            `json:"id"`
		Token            string           `json:"token"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano"`
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
		ExpiresUTCNano   int64            `json:"expires_utc_nano"`
	}

	GetMeResponse struct {
		ID               int64            `json:"id"`
		Login            string           `json:"login"`
	}

	DeleteUserRequest struct {
		Token            string           `json:"token"`
		Login            string           `json:"login"`
		Password         string           `json:"password"`
	}

	DeleteUserResponse struct {
		Description      string           `json:"description"`
	}

	AuthResponse struct {
		Expired          bool             `json:"expired"`
		User             *User            `json:"user.omitempty"`
	}
)