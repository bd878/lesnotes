package middleware

import (
	"net/http"
)

type middleware struct {
	handler  Handler
	funcs    []MiddlewareFunc
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var first Handler = m.handler
	for i := len(m.funcs) - 1; i >= 0; i-- {
		first = m.funcs[i](first)
	}
	first.Handle(w, req)
}
