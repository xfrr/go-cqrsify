package messaging

import (
	"context"
	"errors"
)

// Query represents an action or intent to change the state of the system.
type Query interface {
	Message

	Reply(ctx context.Context, response Message) error
}

// BaseQuery provides a basic implementation of the Query interface.
type BaseQuery struct {
	BaseMessage

	replyCh chan<- Message
}

func (q BaseQuery) Reply(ctx context.Context, response Message) error {
	if q.replyCh == nil {
		return errors.New("no reply channel available")
	}
	select {
	case q.replyCh <- response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// NewBaseQuery creates a new BaseQuery with the given name and payload.
func NewBaseQuery(queryType string, modifiers ...BaseMessageModifier) BaseQuery {
	return BaseQuery{
		BaseMessage: NewBaseMessage(
			queryType,
			modifiers...,
		),
	}
}
