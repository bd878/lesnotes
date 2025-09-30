package model

import "encoding/json"

type (
	ErrorCode struct {
		Code                int              `json:"code"`
		Explain             string           `json:"explain"`
		Human               string           `json:"human"`
	}

	ServerResponse struct {
		Status              string           `json:"status"`
		Error               *ErrorCode       `json:"error,omitempty"`
		Response            json.RawMessage  `json:"response,omitempty"`
		Data                json.RawMessage  `json:"data,omitempty"`
	}

	ServerRequest struct {
		Token               string           `json:"token"`
		Request             json.RawMessage  `json:"req"`
	}
)

