package http

import (
	"net/http"
)

func (h *Handler) SearchFiles(w http.ResponseWriter, req *http.Request) (err error) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
	return nil
}