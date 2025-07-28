package http

import (
        "net/http"
        "time"

        "github.com/bd878/gallery/server/logger"
        httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
)

type logHandler struct {
        next httpmiddleware.Handler
}

var _ httpmiddleware.Handler = (*logHandler)(nil)

func LogBuilder() httpmiddleware.MiddlewareFunc {
        return func(next httpmiddleware.Handler) httpmiddleware.Handler {
                return logHandler{next: next}
        }
}

func (l logHandler) Handle(w http.ResponseWriter, req *http.Request) (err error) {
        logger.Infow("-->", "url", req.URL.String(), "time", time.Now().UnixMilli())
        err = l.next.Handle(w, req)
        logger.Infow("<--", "url", req.URL.String(), "time", time.Now().UnixMilli())
        if err != nil {
                logger.Errorln(err.Error())
        }
        return
}
