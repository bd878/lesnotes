package middleware

import (
	"net/http"
)

type Handler interface {
	Handle(w http.ResponseWriter, req *http.Request) error
}
type MiddlewareFunc func(next Handler) Handler

type HandleFunc func(w http.ResponseWriter, req *http.Request) error 
type handler HandleFunc
func (h handler) Handle(w http.ResponseWriter, req *http.Request) (err error) {
	return h(w, req)
}

type builder struct {
	funcs *list
}

func NewBuilder() builder {
	return builder{
		funcs: &list{},
	}
}

func (b builder) Build(h HandleFunc) *middleware {
	funcs := make([]MiddlewareFunc, 0)
	b.funcs.Traverse(func (n *node) bool {
		if n.f != nil {
			funcs = append(funcs, n.f)
		}
		return false
	})

	return &middleware{
		handler: handler(h),
		funcs: funcs,
	}
}

// Removes auth from all subsequent builds
func (b builder) NoAuth() builder {
	b.funcs.Traverse(func(n *node) bool {
		if n.name == "auth" {
			n.f = nil
		}
		return false
	})

	return b
}

func (b builder) NoLog() builder {
	b.funcs.Traverse(func(n *node) bool {
		if n.name == "log" {
			n.f = nil
		}
		return false
	})

	return b
}

func (b builder) WithLog(f MiddlewareFunc) builder {
	b.funcs.Append(&node{
		name: "log",
		f: f,
	})

	return b
}

func (b builder) WithLang(f MiddlewareFunc) builder {
	b.funcs.Append(&node{
		name: "lang",
		f: f,
	})

	return b
}

func (b builder) WithAuth(f MiddlewareFunc) builder {
	b.funcs.Append(&node{
		name: "auth",
		f: f,
	})

	return b
}