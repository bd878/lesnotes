package model

type (
	Session struct {
		UserID           int64   `json:"user_id"`
		Token            string  `json:"token"`
		ExpiresUTCNano   int64   `json:"expires_utc_nano"`
	}
)