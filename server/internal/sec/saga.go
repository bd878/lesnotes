package sec

type (
	SagaContext[T any] struct {
		ID string
		Data T
		Step int
		Done bool
		Compensating bool
	}
)