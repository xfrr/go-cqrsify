package messaginghttp

import (
	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

// CreateBaseQueryFromSingleDocument creates a BaseQuery from a JSON:API single document.
// It extracts the ID from the document and sets the source to "cqrsify.http".
//
// - If the document contains a "source" meta field, it will be used as the source instead.
//
// - If the document contains a "timestamp" meta field, it will be parsed and used as the timestamp instead of time.Now().
//
// - The "timestamp" meta field must be in RFC3339 format.
//
// Any other meta fields in the document or resource will be added to the query's metadata.
//
// - Meta fields "schema", "source", and "timestamp" are reserved and will not be included in the metadata map.
//
// - If a Meta field is defined both in the resource and the document, the document's value takes precedence.
func CreateBaseQueryFromSingleDocument[A any](queryType string, sd apix.SingleDocument[A]) messaging.BaseQuery {
	schema, source, timestamp, metadata := extractJSONAPISingleDocumentMetadata(sd)

	return messaging.NewBaseQuery(
		queryType,
		messaging.WithQueryID(sd.Data.ID),
		messaging.WithMetadata(metadata),
		messaging.WithSchema(schema),
		messaging.WithSource(source),
		messaging.WithTimestamp(timestamp),
	)
}
