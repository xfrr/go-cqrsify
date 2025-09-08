package message_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfrr/go-cqrsify/message"
)

func TestNewBase(t *testing.T) {
	tests := []struct {
		name      string
		modifiers []message.BaseModifier
		validate  func(t *testing.T, base message.Base)
	}{
		{
			name:      "creates base with default values",
			modifiers: nil,
			validate: func(t *testing.T, base message.Base) {
				assert.Empty(t, base.ID())
				assert.Empty(t, base.Schema())
				assert.Empty(t, base.Source())
				assert.NotZero(t, base.Timestamp())
				assert.NotNil(t, base.Metadata())
				assert.Empty(t, base.Metadata())

				// Verify timestamp is recent (within last second)
				timeDiff := time.Since(base.Timestamp())
				assert.True(t, timeDiff >= 0 && timeDiff < time.Second,
					"Expected timestamp to be recent, but got %v", timeDiff)

				// Verify timestamp is UTC
				assert.Equal(t, time.UTC, base.Timestamp().Location())
			},
		},
		{
			name:      "creates base with empty modifiers slice",
			modifiers: []message.BaseModifier{},
			validate: func(t *testing.T, base message.Base) {
				assert.Empty(t, base.ID())
				assert.Empty(t, base.Schema())
				assert.Empty(t, base.Source())
				assert.NotZero(t, base.Timestamp())
				assert.NotNil(t, base.Metadata())
				assert.Empty(t, base.Metadata())
			},
		},
		{
			name: "creates base with single modifier",
			modifiers: []message.BaseModifier{
				message.WithID("test-id"),
			},
			validate: func(t *testing.T, base message.Base) {
				assert.Equal(t, "test-id", base.ID())
				assert.Empty(t, base.Schema())
				assert.Empty(t, base.Source())
				assert.NotZero(t, base.Timestamp())
				assert.NotNil(t, base.Metadata())
				assert.Empty(t, base.Metadata())
			},
		},
		{
			name: "creates base with multiple modifiers",
			modifiers: []message.BaseModifier{
				message.WithID("test-id"),
				message.WithSchema("test-schema"),
				message.WithSource("test-source"),
				message.WithTimestamp(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
				message.WithMetadata(map[string]string{"key1": "value1"}),
				message.WithMetadataKeyValue("key2", "value2"),
			},
			validate: func(t *testing.T, base message.Base) {
				assert.Equal(t, "test-id", base.ID())
				assert.Equal(t, "test-schema", base.Schema())
				assert.Equal(t, "test-source", base.Source())
				assert.Equal(t, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), base.Timestamp())
				assert.NotNil(t, base.Metadata())
				assert.Len(t, base.Metadata(), 2)
				assert.Equal(t, "value1", base.Metadata()["key1"])
				assert.Equal(t, "value2", base.Metadata()["key2"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := message.NewBase(tt.modifiers...)
			tt.validate(t, base)
		})
	}
}

func TestBase_GettersConsistency(t *testing.T) {
	t.Run("multiple calls to getters return same values", func(t *testing.T) {
		base := message.NewBase()

		// Call getters multiple times and verify consistency
		id1, id2 := base.ID(), base.ID()
		schema1, schema2 := base.Schema(), base.Schema()
		source1, source2 := base.Source(), base.Source()
		timestamp1, timestamp2 := base.Timestamp(), base.Timestamp()
		metadata1, metadata2 := base.Metadata(), base.Metadata()

		assert.Equal(t, id1, id2)
		assert.Equal(t, schema1, schema2)
		assert.Equal(t, source1, source2)
		assert.Equal(t, timestamp1, timestamp2)
		assert.Equal(t, metadata1, metadata2)
	})
}
