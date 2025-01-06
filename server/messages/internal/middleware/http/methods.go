package middleware

import (
  "context"
  "net/http"
  "github.com/bd878/gallery/server/logger"
)

func Log(next Handler) Handler {
  return func(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
    next(ctx, log, w, req)
  }
}

func Auth(next Handler) Handler {
  return func(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
    next(ctx, log, w, req)
  }
}
