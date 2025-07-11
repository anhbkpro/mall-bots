package sec

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type (
	StepActionFunc[T any]       func(ctx context.Context, data T) am.Command
	StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error

	// SagaStep is a step in a saga.
	SagaStep[T any] interface {
		// Action is used to add an action to the step.
		Action(fn StepActionFunc[T]) SagaStep[T]
		// Compensation is used to add a compensation to the step.
		Compensation(fn StepActionFunc[T]) SagaStep[T]
		// OnActionReply is used to add a reply handler to the step.
		OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		// OnCompensationReply is used to add a compensation reply handler to the step.
		OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		// isInvocable is used to check if the step is invocable.
		isInvocable(compensating bool) bool
		// execute is used to execute the step.
		execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T]
		// handle is used to handle a reply from the step.
		handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error
	}

	sagaStep[T any] struct {
		actions  map[bool]StepActionFunc[T]
		handlers map[bool]map[string]StepReplyHandlerFunc[T]
	}

	stepResult[T any] struct {
		ctx *SagaContext[T]
		cmd am.Command
		err error
	}
)

var _ SagaStep[any] = (*sagaStep[any])(nil)

func (s *sagaStep[T]) Action(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep[T]) Compensation(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep[T]) OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s sagaStep[T]) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s sagaStep[T]) execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T] {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		return stepResult[T]{
			ctx: sagaCtx,
			cmd: action(ctx, sagaCtx.Data),
		}
	}

	return stepResult[T]{ctx: sagaCtx}
}

func (s sagaStep[T]) handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error {
	if handler := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaCtx.Data, reply)
	}
	return nil
}

type StepOption[T any] func(step *sagaStep[T])

func WithAction[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[notCompensating] = fn
	}
}

func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[isCompensating] = fn
	}
}

func OnActionReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[notCompensating][replyName] = fn
	}
}

func OnCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[isCompensating][replyName] = fn
	}
}
