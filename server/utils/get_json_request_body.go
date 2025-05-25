package utils

import (
	"net/http"
	"encoding/json"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
)

func GetJsonRequestBody(w http.ResponseWriter, req *http.Request) (json.RawMessage, bool) {
	jsonReq, ok := req.Context().Value(httpmiddleware.RequestContextKey{}).(json.RawMessage)
	if !ok {
		return nil, false
	}
	return jsonReq, true
}