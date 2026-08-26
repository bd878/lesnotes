package sec

import (
	"context"
	"errors"
	"log/slog"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"
)

type (
	Orchestrator interface {
		Start(ctx context.Context, id string, data []byte) error
		ReplyTopic() string
		HandleReply(ctx context.Context, reply ddd.Reply) error
	}

	orchestrator struct {
		saga Saga
		repo SagaRepository
		publisher am.CommandPublisher
	}
)

var _ Orchestrator = (*orchestrator)(nil)

func NewOrchestrator(saga Saga, repo SagaRepository, publisher am.CommandPublisher) Orchestrator {
	return orchestrator{
		saga:      saga,
		repo:      repo,
		publisher: publisher,
	}
}

func (o orchestrator) Start(ctx context.Context, id string, data []byte) error {
	sagaCtx := &SagaContext{
		ID: id,
		Data: data,
		Step: -1,
	}

	err := o.repo.Save(ctx, o.saga.Name(), sagaCtx)
	if err != nil {
		return err
	}

	result := o.execute(ctx, sagaCtx)
	if result.err != nil {
		return err
	}

	return o.processResult(ctx, result)
}

func (o orchestrator) ReplyTopic() string {
	return o.saga.ReplyTopic()
}

func (o orchestrator) HandleReply(ctx context.Context, reply ddd.Reply) error {
	sagaID, sagaName := o.getSagaInfoFromReply(reply)
	slog.Debug("handle reply", slog.String("sagaID", sagaID), slog.String("sagaName", sagaName))
	if sagaID == "" || sagaName == "" || sagaName != o.saga.Name() {
		return nil
	}

	sagaCtx, err := o.repo.Load(ctx, o.saga.Name(), sagaID)
	if err != nil {
		return err
	}

	result, err := o.handle(ctx, sagaCtx, reply)
	if err != nil {
		return err
	}

	return o.processResult(ctx, result)
}

func (o orchestrator) handle(ctx context.Context, sagaCtx *SagaContext, reply ddd.Reply) (stepResult, error) {
	slog.Debug("handle reply", slog.String("name", reply.ReplyName()))
	step := o.saga.getSteps()[sagaCtx.Step]

	err := step.handle(ctx, sagaCtx, reply)
	if err != nil {
		slog.Error("step failed to handle", slog.String("error", err.Error()))
		return stepResult{}, err
	}

	var success bool
	if outcome, ok := reply.Metadata().Get(am.ReplyOutcomeHdr).(string); !ok {
		success = false
	} else {
		success = outcome == am.OutcomeSuccess
	}

	switch {
	case success:
		return o.execute(ctx, sagaCtx), nil
	case sagaCtx.Compensating:
		return stepResult{}, errors.New("received failed reply but already compensating")
	default:
		sagaCtx.compensate()
		return o.execute(ctx, sagaCtx), nil
	}
}

func (o orchestrator) execute(ctx context.Context, sagaCtx *SagaContext) stepResult {
	var delta = 1
	var direction = 1
	var step SagaStep

	if sagaCtx.Compensating {
		direction = -1
	}

	steps := o.saga.getSteps()
	stepCount := len(steps)

	slog.Debug("step count", slog.Int("count", stepCount))

	for i := sagaCtx.Step + direction; i > -1 && i < stepCount; i += direction {
		if step = steps[i]; step != nil && step.isInvocable(sagaCtx.Compensating) {
			break
		}
		delta += 1
	}


	if step == nil {
		slog.Debug("no step, complete")
		sagaCtx.complete()
		return stepResult{ctx: sagaCtx}
	}

	slog.Debug("advance", slog.Int("delta", delta))

	sagaCtx.advance(delta)

	return step.execute(ctx, sagaCtx)
}

func (o orchestrator) processResult(ctx context.Context, result stepResult) (err error) {
	// last step return no command
	if result.cmd != nil {
		err = o.publishCommand(ctx, result)
		if err != nil {
			return
		}
	}

	return o.repo.Save(ctx, o.saga.Name(), result.ctx)
}

func (o orchestrator) publishCommand(ctx context.Context, result stepResult) error {
	cmd := result.cmd

	cmd.Metadata().Set(am.CommandReplyChannelHdr, o.saga.ReplyTopic())
	cmd.Metadata().Set(SagaCommandIDHdr, result.ctx.ID)
	cmd.Metadata().Set(SagaCommandNameHdr, o.saga.Name())

	return o.publisher.Publish(ctx, cmd.Destination(), cmd)
}

func (o orchestrator) getSagaInfoFromReply(reply ddd.Reply) (string, string) {
	var ok bool
	var sagaID, sagaName string

	if sagaID, ok = reply.Metadata().Get(SagaReplyIDHdr).(string); !ok {
		return "", ""
	}

	if sagaName, ok = reply.Metadata().Get(SagaReplyNameHdr).(string); !ok {
		return "", ""
	}

	return sagaID, sagaName
}
