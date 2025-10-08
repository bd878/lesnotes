package ddd

import (
	"time"
	"github.com/google/uuid"
)

type (
	EventPayload interface {}

	Event interface {
		IDer
		EventName()  string
		Payload()    EventPayload
		OccurredAt() time.Time
	}

	IDer interface {
		ID() string
	}

	event struct {
		id         string
		name       string
		occurredAt time.Time
		payload    EventPayload
	}
)

var _ Event = (*event)(nil)

func NewEvent(name string, payload EventPayload) event {
	return event{
		id:         uuid.New().String(),
		name:       name,
		occurredAt: time.Now(),
		payload:    payload,
	}
}

func (e event) ID() string { return e.id }
func (e event) EventName() string { return e.name }
func (e event) OccurredAt() time.Time { return e.occurredAt }
func (e event) Payload() EventPayload { return e.payload }
