package model

import "encoding/json"

type ServerResponse struct {
	Status              string        `json:"status"`
	Description         string        `json:"description"`
}

type JSONServerRequest struct {
	Token               string          `json:"token"`
	Req                 json.RawMessage `json:"req"`
}