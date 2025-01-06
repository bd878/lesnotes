package middleware

type node struct {
  name  string
  next *node
  f     MiddlewareFunc
}

type list struct {
  len     int
  head   *node
  last   *node /* append only */
}

func (l *list) Len() int {
  return l.len
}

func (l *list) Traverse(fn func(n *node) bool) (found bool) {
  for n := l.head; n != nil; n = n.next {
    if fn(n) {
      found = true
      return
    }
  }
  found = false
  return
}

func (l *list) Append(n *node) {
  l.len += 1

  if l.head == nil {
    l.head = n
    l.last = l.head
  } else {
    l.last.next = n
    l.last = n
  }
}
