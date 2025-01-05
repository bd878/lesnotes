package middleware

import (
  "context"
  "net/http"
  "github.com/bd878/gallery/server/logger"
)

func Logging(
  next func (ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request),
) (
  func (w http.ResponseWriter, req *http.Request),
) {
  return func(w http.ResponseWriter, req *http.Request) {
    next(context.Background(), logger.Default(), w, req)
  }
}
