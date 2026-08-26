package sec

import (
	"context"
)

type SagaStore interface {
	Load(ctx context.Context, sagaName, sagaID string) (*SagaContext, error)
	Save(ctx context.Context, sagaName string, sagaCtx *SagaContext) error
}

type SagaRepository struct {
	store SagaStore
}

func NewSagaRepository(store SagaStore) SagaRepository {
	return SagaRepository{
		store: store,
	}
}

func (r SagaRepository) Load(ctx context.Context, sagaName, sagaID string) (*SagaContext, error) {
	byteCtx, err := r.store.Load(ctx, sagaName, sagaID)
	if err != nil {
		return nil, err
	}

	return &SagaContext{
		ID:           byteCtx.ID,
		Data:         byteCtx.Data,
		Step:         byteCtx.Step,
		Done:         byteCtx.Done,
		Compensating: byteCtx.Compensating,
	}, nil
}

func (r SagaRepository) Save(ctx context.Context, sagaName string, sagaCtx *SagaContext) error {
	return r.store.Save(ctx, sagaName, &SagaContext{
		ID:           sagaCtx.ID,
		Data:         sagaCtx.Data,
		Step:         sagaCtx.Step,
		Done:         sagaCtx.Done,
		Compensating: sagaCtx.Compensating,
	})
}
