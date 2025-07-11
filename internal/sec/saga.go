package sec

import (
	"eda-in-golang/internal/am"
	"fmt"
)

const (
	SagaCommandIDHdr   = am.CommandHdrPrefix + "SAGA_ID"
	SagaCommandNameHdr = am.CommandHdrPrefix + "SAGA_NAME"

	SagaReplyIDHdr   = am.ReplyHdrPrefix + "SAGA_ID"
	SagaReplyNameHdr = am.ReplyHdrPrefix + "SAGA_NAME"
)

type (
	SagaContext[T any] struct {
		ID           string
		Data         T
		Step         int
		Done         bool
		Compensating bool
	}

	Saga[T any] interface {
		AddStep() SagaStep[T]
		Name() string
		ReplyTopic() string
		getSteps() []SagaStep[T]
	}

	saga[T any] struct {
		name       string
		replyTopic string
		steps      []SagaStep[T]
	}
)

const (
	notCompensating = false
	isCompensating  = true
)

func NewSaga[T any](name, replyTopic string) Saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

func (s *saga[T]) AddStep() SagaStep[T] {
	step := &sagaStep[T]{
		actions: map[bool]StepActionFunc[T]{
			notCompensating: nil,
			isCompensating:  nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc[T]{
			notCompensating: {},
			isCompensating:  {},
		},
	}

	s.steps = append(s.steps, step)

	return step
}

func (s *saga[T]) Name() string {
	fmt.Println("=== [Sec] Getting saga name")
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	fmt.Println("=== [Sec] Getting reply topic")
	return s.replyTopic
}

func (s *saga[T]) getSteps() []SagaStep[T] {
	fmt.Println("=== [Sec] Getting steps")
	return s.steps
}

func (s *SagaContext[T]) advance(steps int) {
	fmt.Printf("=== [Sec] Advancing saga from step %d by %d steps (compensating: %v)\n", s.Step, steps, s.Compensating)
	var dir = 1
	if s.Compensating {
		dir = -1
	}

	s.Step += dir * steps
	fmt.Printf("=== [Sec] Advanced to step %d\n", s.Step)
}

func (s *SagaContext[T]) complete() {
	fmt.Println("=== [Sec] Completing saga")
	s.Done = true
}

func (s *SagaContext[T]) compensate() {
	fmt.Println("=== [Sec] Compensating saga")
	s.Compensating = true
}
