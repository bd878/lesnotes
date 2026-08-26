package sec

import "github.com/bd878/gallery/server/internal/am"

const (
	SagaCommandIDHdr   = am.CommandHdrPrefix + "SAGA_ID"
	SagaCommandNameHdr = am.CommandHdrPrefix + "SAGA_NAME"

	SagaReplyIDHdr   = am.ReplyHdrPrefix + "SAGA_ID"
	SagaReplyNameHdr = am.ReplyHdrPrefix + "SAGA_NAME"
)

type (
	SagaContext struct {
		ID string
		Data []byte
		Step int
		Done bool
		Compensating bool
	}

	Saga interface {
		AddStep() SagaStep
		Name() string
		ReplyTopic() string
		getSteps() []SagaStep
	}

	saga struct {
		name string
		replyTopic string
		steps []SagaStep
	}
)

const (
	notCompensating = false
	isCompensating = true
)

func NewSaga(name, replyTopic string) Saga {
	return &saga{
		name: name,
		replyTopic: replyTopic,
	}
}

func (s *saga) AddStep() SagaStep {
	step := &sagaStep{
		actions: map[bool]StepActionFunc{
			notCompensating: nil,
			isCompensating: nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc{
			notCompensating: {},
			isCompensating: {},
		},
	}

	s.steps = append(s.steps, step)

	return step
}

func (s *saga) Name() string {
	return s.name
}

func (s *saga) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga) getSteps() []SagaStep {
	return s.steps
}

func (s *SagaContext) advance(steps int) {
	var direction = 1
	if s.Compensating {
		direction = -1
	}

	s.Step += direction * steps
}

func (s *SagaContext) complete() {
	s.Done = true
}

func (s *SagaContext) compensate() {
	s.Compensating = true
}