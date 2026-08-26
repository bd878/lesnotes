package sec

import (
	"context"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	StepActionFunc func(ctx context.Context, data []byte) am.Command
	StepReplyHandlerFunc func(ctx context.Context, data []byte, reply ddd.Reply) error

	SagaStep interface {
		Action(fn StepActionFunc) SagaStep
		Compensation(fn StepActionFunc) SagaStep
		OnActionReply(replyName string, fn StepReplyHandlerFunc) SagaStep
		OnCompensationReply(replyName string, fn StepReplyHandlerFunc) SagaStep
		isInvocable(compensating bool) bool
		execute(ctx context.Context, sagaCtx *SagaContext) stepResult
		handle(ctx context.Context, sagaCtx *SagaContext, reply ddd.Reply) error
	}

	sagaStep struct {
		actions map[bool]StepActionFunc
		handlers map[bool]map[string]StepReplyHandlerFunc
	}

	stepResult struct {
		ctx *SagaContext
		cmd am.Command
		err error
	}
)

var _ SagaStep = (*sagaStep)(nil)

func (s *sagaStep) Action(fn StepActionFunc) SagaStep {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep) Compensation(fn StepActionFunc) SagaStep {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep) OnActionReply(replyName string, fn StepReplyHandlerFunc) SagaStep {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s *sagaStep) OnCompensationReply(replyName string, fn StepReplyHandlerFunc) SagaStep {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s sagaStep) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s sagaStep) execute(ctx context.Context, sagaCtx *SagaContext) stepResult {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		return stepResult{
			ctx: sagaCtx,
			cmd: action(ctx, sagaCtx.Data),
		}
	}

	return stepResult{ctx: sagaCtx}
}

func (s sagaStep) handle(ctx context.Context, sagaCtx *SagaContext, reply ddd.Reply) error {
	if handler := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaCtx.Data, reply)
	}
	return nil
}

type StepOption func(step *sagaStep)

func WithAction(fn StepActionFunc) StepOption {
	return func(step *sagaStep) {
		step.actions[notCompensating] = fn
	}
}

func WithCompensation(fn StepActionFunc) StepOption {
	return func(step *sagaStep) {
		step.actions[isCompensating] = fn
	}
}

func OnActionReply(replyName string, fn StepReplyHandlerFunc) StepOption {
	return func(step *sagaStep) {
		step.handlers[notCompensating][replyName] = fn
	}
}

func OnCompensationReply(replyName string, fn StepReplyHandlerFunc) StepOption {
	return func(step *sagaStep) {
		step.handlers[isCompensating][replyName] = fn
	}
}
