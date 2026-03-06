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

func (s *Reader) ReadSize() (size uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buf := make([]byte, lenWidth)
	if _, err = s.ReadCloser.Read(buf); err != nil {
		return 0, err
	}
	size = enc.Uint64(buf)
	return
}

func (s *Reader) Read(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n, err = s.ReadCloser.Read(p); err != nil {
		return 0, err
	}
	return
}

func (s *Reader) Close() (err error) {
	return s.ReadCloser.Close()
}