package model

const (
	PublicUserID int32 = 9_999_999
)

type (
	User struct {
		ID               int32            `json:"id"`
		Name             string           `json:"name,omitempty"`
		Password         string           `json:"password,omitempty"`
		Token            string           `json:"token,omitempty"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano,omitempty"`
	}

	ServerResponse struct {
		Status           string           `json:"status"`
		Description      string           `json:"description"`
	}

	ServerAuthorizeResponse struct {
		ServerResponse
		Expired          bool             `json:"expired"`
		User             User             `json:"user"`
	}
)