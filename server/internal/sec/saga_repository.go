package sec

import (
	"context"
)

type SagaStore interface {
	Load(ctx context.Context, sagaName, sagaID string) (*SagaContext[[]byte], error)
	Save(ctx context.Context, sagaName string, sagaCtx *SagaContext[[]byte]) error
}