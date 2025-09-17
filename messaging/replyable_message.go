package messaging

import "context"

const DefaultReplyTimeoutSeconds = 5

// ReplyableMessage represents a message that can receive replies.
type ReplyableMessage interface {
	Message

	// GetReply returns the reply channel associated with the message.
	GetReply(ctx context.Context) (Message, error)
}
