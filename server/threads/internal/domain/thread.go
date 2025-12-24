package domain

import (
	"errors"
	"github.com/bd878/gallery/server/ddd"
)

const (
	ThreadCreatedEvent = "threads.ThreadCreated"
	ThreadReorderedEvent = "threads.ThreadReordered"
)

var (
	ErrIDRequired = errors.New("id is 0")
)

type ThreadCreated struct {

}

func (ThreadCreated) Key() string { return ThreadCreatedEvent }

func CreateThread(id int64) {}

type ThreadReordered struct {

}

func (ThreadReordered) Key() string { return ThreadReorderedEvent }

func ReorderThread() (ddd.Event, error) {
	
}