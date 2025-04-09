package model

import "github.com/bd878/gallery/server/pkg/model"

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

	ServerAuthorizeResponse struct {
		model.ServerResponse
		Expired          bool             `json:"expired"`
		User             User             `json:"user"`
	}

	ServerUserResponse struct {
		model.ServerResponse
		User             User             `json:"user"`
	}
)