package ddd

import (
	"time"
	"github.com/google/uuid"
)

type (
	EventOption interface {
		configureEvent(*event)
	}

	EventPayload interface {}

	Event interface {
		IDer
		EventName()  string
		Payload()    EventPayload
		Metadata()   Metadata
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
		metadata   Metadata
	}
)

var _ Event = (*event)(nil)

func NewEvent(name string, payload EventPayload, options ...EventOption) event {
	evt := event{
		id:         uuid.New().String(),
		name:       name,
		occurredAt: time.Now(),
		payload:    payload,
		metadata:   make(Metadata),
	}

	for _, option := range options {
		option.configureEvent(&evt)
	}

	return evt
}

func (e event) ID() string { return e.id }
func (e event) EventName() string { return e.name }
func (e event) OccurredAt() time.Time { return e.occurredAt }
func (e event) Payload() EventPayload { return e.payload }
func (e event) Metadata() Metadata { return e.metadata }
