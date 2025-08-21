package model

import "github.com/bd878/gallery/server/pkg/model"

const (
	PublicUserID int64 = 9_999_999
)

type (
	User struct {
		ID               int64            `json:"id"`
		Name             string           `json:"name,omitempty"` // TODO: login
		Password         string           `json:"password,omitempty"` // TODO: HashedPassword
		Token            string           `json:"token,omitempty"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano,omitempty"` // TODO: ExpiresAt
	}

	DeleteUserJsonRequest struct {
		Token            string           `json:"token"`
		Name             string           `json:"name"`
		Password         string           `json:"password"`
	}

	DeleteUserJsonServerResponse struct {
		model.ServerResponse
		Expired          bool             `json:"expired,omitempty"`
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

	SignupJsonUserServerResponse struct {
		model.ServerResponse
		ID               int64            `json:"id"`
		Token            string           `json:"token"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano"`
	}

	LoginJsonUserServerResponse struct {
		model.ServerResponse
		Token            string           `json:"token"`
		ExpiresUTCNano   int64            `json:"expires_utc_nano"`
	}
)