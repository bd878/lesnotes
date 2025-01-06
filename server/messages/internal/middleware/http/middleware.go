package middleware

import (
  "net/http"
  "context"
  "github.com/bd878/gallery/server/logger"
)

type Middleware struct {
  handler  Handler
  funcs    []MiddlewareFunc
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  var first Handler = m.handler
  for i := len(m.funcs) - 1; i >= 0; i-- {
    first = m.funcs[i](first)
  }
  first(context.Background(), logger.Default(), w, req)
}
