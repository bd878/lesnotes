package middleware

import (
	"net/http"

	"github.com/bd878/gallery/server/logger"
)

func Log(next Handler) Handler {
	return handler(func(log *logger.Logger, w http.ResponseWriter, req *http.Request) (err error) {
		log.Infof("--> %s\n", req.URL.String())
		err = next.Handle(log, w, req)
		log.Infof("<-- %s\n", req.URL.String())
		if err != nil {
			log.Errorln(err.Error())
		}
		return
	})
}
