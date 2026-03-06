package store

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8 // 8 bytes
)

type Store struct {
	*os.File
	mu     sync.Mutex
	buf    *bufio.Writer
	pos    uint64
}

func NewStore(f *os.File) (s *Store, err error) {
	return &Store{
		File:     f,
		buf:      bufio.NewWriter(f),
	}, nil
}

func (s *Store) Append(p []byte) (n uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, err
	}
	w += lenWidth
	return uint64(w), nil
}

func (s *Store) Read(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	size := make([]byte, lenWidth)
	if n, err = s.File.ReadAt(size, int64(s.pos)); err != nil {
		return 0, err
	}
	p = make([]byte, enc.Uint64(size))
	if n, err = s.File.ReadAt(p, int64(s.pos + lenWidth)); err != nil {
		return 0, err
	}
	s.pos += uint64(lenWidth + n)
	return
}

func (s *Store) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err = s.buf.Flush()
	if err != nil {
		return
	}
	return s.File.Close()
}