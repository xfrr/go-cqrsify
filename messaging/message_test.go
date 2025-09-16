package messaging_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestNewBaseMessage(t *testing.T) {
	tests := []struct {
		name      string
		modifiers []messaging.BaseMessageModifier
		validate  func(t *testing.T, base messaging.BaseMessage)
	}{
		{
			name:      "creates base with default values",
			modifiers: nil,
			validate: func(t *testing.T, base messaging.BaseMessage) {
				assert.Empty(t, base.MessageID())
				assert.Empty(t, base.MessageSchemaURI())
				assert.Empty(t, base.MessageSource())
				assert.NotZero(t, base.MessageTimestamp())
				assert.NotNil(t, base.MessageMetadata())
				assert.Empty(t, base.MessageMetadata())

				// Verify timestamp is recent (within last second)
				timeDiff := time.Since(base.MessageTimestamp())
				assert.True(t, timeDiff >= 0 && timeDiff < time.Second,
					"Expected timestamp to be recent, but got %v", timeDiff)

				// Verify timestamp is UTC
				assert.Equal(t, time.UTC, base.MessageTimestamp().Location())
			},
		},
		{
			name:      "creates base with empty modifiers slice",
			modifiers: []messaging.BaseMessageModifier{},
			validate: func(t *testing.T, base messaging.BaseMessage) {
				assert.Empty(t, base.MessageID())
				assert.Empty(t, base.MessageSchemaURI())
				assert.Empty(t, base.MessageSource())
				assert.NotZero(t, base.MessageTimestamp())
				assert.NotNil(t, base.MessageMetadata())
				assert.Empty(t, base.MessageMetadata())
			},
		},
		{
			name: "creates base with single modifier",
			modifiers: []messaging.BaseMessageModifier{
				messaging.WithID("test-id"),
			},
			validate: func(t *testing.T, base messaging.BaseMessage) {
				assert.Equal(t, "test-id", base.MessageID())
				assert.Empty(t, base.MessageSchemaURI())
				assert.Empty(t, base.MessageSource())
				assert.NotZero(t, base.MessageTimestamp())
				assert.NotNil(t, base.MessageMetadata())
				assert.Empty(t, base.MessageMetadata())
			},
		},
		{
			name: "creates base with multiple modifiers",
			modifiers: []messaging.BaseMessageModifier{
				messaging.WithID("test-id"),
				messaging.WithSchema("test-schema"),
				messaging.WithSource("test-source"),
				messaging.WithTimestamp(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
				messaging.WithMetadata(map[string]string{"key1": "value1"}),
				messaging.WithMetadataKeyValue("key2", "value2"),
			},
			validate: func(t *testing.T, base messaging.BaseMessage) {
				assert.Equal(t, "test-id", base.MessageID())
				assert.Equal(t, "test-schema", base.MessageSchemaURI())
				assert.Equal(t, "test-source", base.MessageSource())
				assert.Equal(t, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), base.MessageTimestamp())
				assert.NotNil(t, base.MessageMetadata())
				assert.Len(t, base.MessageMetadata(), 2)
				assert.Equal(t, "value1", base.MessageMetadata()["key1"])
				assert.Equal(t, "value2", base.MessageMetadata()["key2"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := messaging.NewBaseMessage(tt.name, tt.modifiers...)
			tt.validate(t, base)
		})
	}
}

func TestBaseMessage_GettersConsistency(t *testing.T) {
	t.Run("multiple calls to getters return same values", func(t *testing.T) {
		base := messaging.NewBaseMessage("test-type")

		// Call getters multiple times and verify consistency
		id1, id2 := base.MessageID(), base.MessageID()
		schema1, schema2 := base.MessageSchemaURI(), base.MessageSchemaURI()
		source1, source2 := base.MessageSource(), base.MessageSource()
		timestamp1, timestamp2 := base.MessageTimestamp(), base.MessageTimestamp()
		metadata1, metadata2 := base.MessageMetadata(), base.MessageMetadata()

		assert.Equal(t, id1, id2)
		assert.Equal(t, schema1, schema2)
		assert.Equal(t, source1, source2)
		assert.Equal(t, timestamp1, timestamp2)
		assert.Equal(t, metadata1, metadata2)
	})
}
