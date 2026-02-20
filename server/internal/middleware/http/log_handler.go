package middleware

import (
	"net/http"

	"github.com/bd878/gallery/server/internal/logger"
)

func Log(next Handler) Handler {
	return handler(func(w http.ResponseWriter, req *http.Request) (err error) {
		logger.Default().Infof("--> %s\n", req.URL.String())
		err = next.Handle(w, req)
		logger.Default().Infof("<-- %s\n", req.URL.String())
		if err != nil {
			logger.Default().Errorln(err.Error())
		}
		return
	})
}
