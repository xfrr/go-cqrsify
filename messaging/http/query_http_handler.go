package messaginghttp

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

// QueryHandler is an alias to HTTPMessageServer to keep external API surface familiar.
type QueryHandler = MessageWithReplyHandler

// NewQueryHandler creates a new QueryHTTPServer with the given QueryDispatcher and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewQueryHandler(dispatcher messaging.QueryDispatcher, opts ...MessageHandlerWithReplyOption) *QueryHandler {
	return NewMessageWithReplyHandler(&queryDispatcherWrapper{dispatcher}, opts...)
}

// RegisterJSONAPIQueryDecoder registers a JSON:API query decoder for the given query type.
// If a decoder for the same query type and encoding already exists, an error is returned.
func RegisterJSONAPIQueryDecoder[A any](handler *QueryHandler, msgType string, decodeFunc func(context.Context, apix.SingleDocument[A]) (messaging.Query, error)) error {
	return RegisterJSONSingleDocumentMessageDecoder(&handler.inner, msgType, func(ctx context.Context, sd apix.SingleDocument[A]) (messaging.Message, error) {
		return decodeFunc(ctx, sd)
	})
}
