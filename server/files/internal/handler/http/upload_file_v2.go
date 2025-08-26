package http

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"

	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) UploadFileV2(w http.ResponseWriter, req *http.Request) (err error) {
	var public int
	values := req.URL.Query()
	if values.Has("public") {
		public, err = strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongFormat,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "public"),
				},
			})

			return err
		}
	} else {
		public = -1
	}

	return h.uploadFile(w, req, public)
}