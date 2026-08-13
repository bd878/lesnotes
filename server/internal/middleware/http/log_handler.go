package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Log(next Handler) Handler {
	return handler(func(w http.ResponseWriter, req *http.Request) (err error) {
		slog.Info(fmt.Sprintf("--> %s", req.URL.String()))
		err = next.Handle(w, req)
		slog.Info(fmt.Sprintf("<-- %s", req.URL.String()))
		if err != nil {
			slog.Error(err.Error())
		}
		return
	})
}
