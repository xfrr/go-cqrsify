package messaging

import (
	"context"
	"errors"
	"fmt"
)

// Query represents an action or intent to change the state of the system.
type Query interface {
	Message

	Reply(ctx context.Context, response Message) error
}

// BaseQuery provides a basic implementation of the Query interface.
type BaseQuery struct {
	baseMessage

	replyCh chan Message
}

type BaseQueryModifier = baseMessageModifier

func (q BaseQuery) GetReply(ctx context.Context) (Message, error) {
	select {
	case reply := <-q.replyCh:
		return reply, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("getting reply timed out: %w", ctx.Err())
	}
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
func NewBaseQuery(queryType string, modifiers ...BaseQueryModifier) BaseQuery {
	return BaseQuery{
		replyCh: make(chan Message, 1),
		baseMessage: newBaseMessage(
			queryType,
			modifiers...,
		),
	}
}

// NewQueryFromJSON creates a BaseQuery from a JSONMessage.
func NewQueryFromJSON[P any](jsonMsg JSONMessage[P]) BaseQuery {
	return BaseQuery{
		replyCh: make(chan Message, 1),
		baseMessage: baseMessage{
			id:        jsonMsg.ID,
			_type:     jsonMsg.Type,
			schema:    jsonMsg.SchemaURI,
			source:    jsonMsg.Source,
			timestamp: jsonMsg.Timestamp,
			metadata:  jsonMsg.Metadata,
		},
	}
}
