package middleware

import (
  "context"
  "net/http"
  "github.com/bd878/gallery/server/logger"
)

type Handler func(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request)
type MiddlewareFunc func(next Handler) Handler

type builder struct {
  funcs *list
}

func NewBuilder() builder {
  return builder{
    funcs: &list{},
  }
}

func (b builder) Build(handler Handler) *Middleware {
  funcs := make([]MiddlewareFunc, 0)
  b.funcs.Traverse(func (n *node) bool {
    if n.f != nil {
      funcs = append(funcs, n.f)
    }
    return false
  })

  return &Middleware{
    handler: handler,
    funcs: funcs,
  }
}

func (b builder) NoAuth() builder {
  b.funcs.Traverse(func(n *node) bool {
    if n.name == "auth" {
      n.f = nil
      return true
    }
    return false
  })

  return b
}

func (b builder) NoLog() builder {
  b.funcs.Traverse(func(n *node) bool {
    if n.name == "log" {
      n.f = nil
      return true
    }
    return false
  })

  return b
}

func (b builder) WithLog() builder {
  found := b.funcs.Traverse(func(n *node) bool {
    if n.name == "log" {
      n.f = nil
      return true
    }
    return false
  })

  if !found {
    b.funcs.Append(&node{
      name: "log",
      f: Log,
    })
  }

  return b
}

func (b builder) WithAuth() builder {
  found := b.funcs.Traverse(func(n *node) bool {
    if n.name == "auth" {
      n.f = nil
      return true
    }
    return false
  })

  if !found {
    b.funcs.Append(&node{
      name: "auth",
      f: Auth,
    })
  }

  return b
}