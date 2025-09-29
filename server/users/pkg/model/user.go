package model

const (
	PublicUserID int64 = 9_999_999
)

type (
	User struct {
		ID               int64            `json:"id"`
		Login            string           `json:"login,omitempty"`
		Theme            string           `json:"theme,omitempty"`
		HashedPassword   string           `json:"salt,omitempty"`
		Token            string           `json:"token,omitempty"`
		Lang             string           `json:"language,omitempty"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano,omitempty"` // TODO: ExpiresAt
		FontSize         int32            `json:"font_size,omitempty"`
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
		Theme            string           `json:"theme"`
		Lang             string           `json:"language"`
		FontSize         int32            `json:"font_size"`
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
		Login            string           `json:"login,omitempty"`
		Theme            string           `json:"theme,omitempty"`
		Lang             string           `json:"language,omitempty"`
		FontSize         int32            `json:"font_size,omitempty"`
	}

	UpdateResponse struct {
		Description      string           `json:"description"`
	}
)