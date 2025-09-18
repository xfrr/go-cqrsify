package messaging

import (
	"context"
	"errors"
	"fmt"
)

// Query represents an action or intent to change the state of the system.
type Query interface {
	Message

	Reply(ctx context.Context, reply QueryReply) error
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

func (q BaseQuery) Reply(ctx context.Context, reply QueryReply) error {
	if q.replyCh == nil {
		return errors.New("no reply channel available")
	}

	select {
	case q.replyCh <- reply:
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

// QueryReply represents a reply to a Query.
type QueryReply interface {
	Message
}

// BaseQueryReply provides a basic implementation of the QueryReply interface.
type BaseQueryReply struct {
	baseMessage
}

type BaseQueryReplyModifier = baseMessageModifier

// NewBaseQueryReply creates a new BaseQueryReply with the given name and payload.
func NewBaseQueryReply(query Query, modifiers ...BaseQueryReplyModifier) BaseQueryReply {
	return BaseQueryReply{
		baseMessage: newBaseMessage(
			query.MessageType()+".reply",
			modifiers...,
		),
	}
}

// NewQueryReplyFromJSON creates a BaseQueryReply from a JSONMessage.
func NewQueryReplyFromJSON[P any](jsonMsg JSONMessage[P]) BaseQueryReply {
	return BaseQueryReply{
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
