package model

import "encoding/json"

type ServerResponse struct {
	Status              string        `json:"status"`
	Code                string        `json:"code,omitempty"`
	Description         string        `json:"description,omitempty"`
}

type JSONServerRequest struct {
	Token               string          `json:"token"`
	Req                 json.RawMessage `json:"req"`
}