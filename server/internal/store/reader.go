package store

import (
	"io"
	"sync"
)

type Reader struct {
	io.ReadCloser
	mu sync.Mutex
}

func NewReader(reader io.ReadCloser) *Reader {
	return &Reader{
		ReadCloser: reader,
	}
}

func (s *Reader) Read(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	size := make([]byte, lenWidth)
	if n, err = s.ReadCloser.Read(size); err != nil {
		return 0, err
	}
	p = make([]byte, enc.Uint64(size))
	if n, err = s.ReadCloser.Read(p); err != nil {
		return 0, err
	}
	n += lenWidth
	return
}

func (s *Reader) Close() (err error) {
	return s.ReadCloser.Close()
}