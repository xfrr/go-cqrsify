package messaginghttp

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

const (
	jsonAPISingleDocumentMetadataSchemaKey    = "schema"
	jsonAPISingleDocumentMetadataSourceKey    = "source"
	jsonAPISingleDocumentMetadataTimestampKey = "timestamp"
)

// CreateBaseCommandFromSingleDocument creates a BaseCommand from a JSON:API single document.
// It extracts the ID from the document and sets the source to "cqrsify.http".
//
// - If the document contains a "source" meta field, it will be used as the source instead.
//
// - If the document contains a "timestamp" meta field, it will be parsed and used as the timestamp instead of time.Now().
//
// - The "timestamp" meta field must be in RFC3339 format.
//
// Any other meta fields in the document or resource will be added to the command's metadata.
//
// - Meta fields "schema", "source", and "timestamp" are reserved and will not be included in the metadata map.
//
// - If a Meta field is defined both in the resource and the document, the document's value takes precedence.
func CreateBaseCommandFromSingleDocument[A any](cmdType string, sd apix.SingleDocument[A]) messaging.BaseCommand {
	var (
		source           string
		schema           string
		timestamp        = time.Now()
		resourceMetadata = make(map[string]string)
		keys             = []string{
			jsonAPISingleDocumentMetadataSchemaKey,
			jsonAPISingleDocumentMetadataSourceKey,
			jsonAPISingleDocumentMetadataTimestampKey,
		}
	)

	if sd.Data.Meta != nil {
		if sch, ok := sd.Data.Meta[jsonAPISingleDocumentMetadataSchemaKey].(string); ok {
			schema = sch
		}

		if s, ok := sd.Data.Meta[jsonAPISingleDocumentMetadataSourceKey].(string); ok {
			source = s
		}

		if ts, ok := sd.Data.Meta[jsonAPISingleDocumentMetadataTimestampKey].(string); ok {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				timestamp = t
			}
		}

		resourceMetadata = filterAndParseMetadata(sd.Data.Meta, keys)
		documentMetadata := filterAndParseMetadata(sd.Meta, keys)
		maps.Copy(resourceMetadata, documentMetadata)
	}

	return messaging.BaseCommand{
		BaseMessage: messaging.NewBaseMessage(
			cmdType,
			messaging.WithID(sd.Data.ID),
			messaging.WithMetadata(resourceMetadata),
			messaging.WithSchema(schema),
			messaging.WithSource(source),
			messaging.WithTimestamp(timestamp),
		),
	}
}

func filterAndParseMetadata(sd map[string]any, excludeKeys []string) map[string]string {
	metadata := make(map[string]string)

	for k, v := range sd {
		if slices.Contains(excludeKeys, k) {
			continue
		}

		vStr, ok := v.(string)
		if !ok {
			// 	check if can be converted to string
			vcStr, vcStrOk := parseAnyString(v)
			if !vcStrOk {
				continue
			}

			vStr = vcStr
		}

		metadata[k] = vStr
	}

	return metadata
}

func parseAnyString(v any) (string, bool) {
	switch val := v.(type) {
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), true
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), true
	case float32, float64:
		return fmt.Sprintf("%f", val), true
	case bool:
		return strconv.FormatBool(val), true
	default:
		return "", false
	}
}
