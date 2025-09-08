package message_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/message"
)

// TestPayload represents a sample payload for testing
type TestPayload struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// createTestBase creates an empty Base message instance
func createTestBase() message.Base {
	return message.Base{}
}

// createTestBaseWithProps creates a Base message instance with specified properties
func createTestBaseWithProps(id, schema, source string, metadata map[string]string, timestamp time.Time) message.Base {
	base := message.NewBase(
		message.WithID(id),
		message.WithSchema(schema),
		message.WithSource(source),
		message.WithMetadata(metadata),
		message.WithTimestamp(timestamp),
	)
	return base
}

func TestNewEnvelope(t *testing.T) {
	tests := []struct {
		name     string
		msgBase  message.Base
		payload  TestPayload
		expected func(message.Base, TestPayload) message.Envelope
	}{
		{
			name:    "creates envelope with valid base and payload",
			msgBase: createTestBase(),
			payload: TestPayload{Name: "test", Value: 42},
			expected: func(base message.Base, payload TestPayload) message.Envelope {
				return message.NewEnvelope(base, payload)
			},
		},
		{
			name:    "creates envelope with empty payload",
			msgBase: createTestBase(),
			payload: TestPayload{},
			expected: func(base message.Base, payload TestPayload) message.Envelope {
				return message.NewEnvelope(base, payload)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := message.NewEnvelope(tt.msgBase, tt.payload)
			expected := tt.expected(tt.msgBase, tt.payload)

			assert.Equal(t, expected.Base, result.Base)
			assert.Equal(t, expected.Payload(), result.Payload())
		})
	}
}

func TestEnvelope_MarshalJSON(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	metadata := map[string]string{"key1": "value1", "key2": "value2"}

	tests := []struct {
		name     string
		envelope message.Envelope
		expected string
		wantErr  bool
	}{
		{
			name: "marshals envelope with complete data",
			envelope: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-123",
					"test-schema",
					"test-source",
					metadata,
					fixedTime,
				),
				TestPayload{Name: "test-payload", Value: 100},
			),
			expected: `{"id":"test-id-123","schema":"test-schema","source":"test-source","metadata":{"key1":"value1","key2":"value2"},"payload":{"name":"test-payload","value":100},"timestamp":"2023-12-25T10:30:00Z"}`,
			wantErr:  false,
		},
		{
			name: "marshals envelope with nil metadata",
			envelope: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-456",
					"test-schema-2",
					"test-source-2",
					nil,
					fixedTime,
				),
				TestPayload{Name: "test-payload-2", Value: 200},
			),
			expected: `{"id":"test-id-456","schema":"test-schema-2","source":"test-source-2","metadata":null,"payload":{"name":"test-payload-2","value":200},"timestamp":"2023-12-25T10:30:00Z"}`,
			wantErr:  false,
		},
		{
			name: "marshals envelope with empty payload",
			envelope: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-789",
					"test-schema-3",
					"test-source-3",
					make(map[string]string),
					fixedTime,
				),
				TestPayload{},
			),
			expected: `{"id":"test-id-789","schema":"test-schema-3","source":"test-source-3","metadata":{},"payload":{"name":"","value":0},"timestamp":"2023-12-25T10:30:00Z"}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.envelope.MarshalJSON()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestEnvelope_UnmarshalJSON(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		jsonData string
		expected message.Envelope
		wantErr  bool
	}{
		{
			name:     "unmarshals valid JSON with complete data",
			jsonData: `{"id":"test-id-123","schema":"test-schema","source":"test-source","metadata":{"key1":"value1","key2":"value2"},"payload":{"name":"test-payload","value":100},"timestamp":"2023-12-25T10:30:00Z"}`,
			expected: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-123",
					"test-schema",
					"test-source",
					map[string]string{"key1": "value1", "key2": "value2"},
					fixedTime,
				),
				map[string]interface{}{"name": "test-payload", "value": float64(100)},
			),
			wantErr: false,
		},
		{
			name:     "unmarshals JSON with null metadata",
			jsonData: `{"id":"test-id-456","schema":"test-schema-2","source":"test-source-2","metadata":null,"payload":{"name":"test-payload-2","value":200},"timestamp":"2023-12-25T10:30:00Z"}`,
			expected: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-456",
					"test-schema-2",
					"test-source-2",
					nil,
					fixedTime,
				),
				map[string]interface{}{"name": "test-payload-2", "value": float64(200)}, // JSON numbers are float64
			),
			wantErr: false,
		},
		{
			name:     "unmarshals JSON with empty metadata",
			jsonData: `{"id":"test-id-789","schema":"test-schema-3","source":"test-source-3","metadata":{},"payload":{"name":"","value":0},"timestamp":"2023-12-25T10:30:00Z"}`,
			expected: message.NewEnvelope(
				createTestBaseWithProps(
					"test-id-789",
					"test-schema-3",
					"test-source-3",
					make(map[string]string),
					fixedTime,
				),
				map[string]interface{}{"name": "", "value": float64(0)},
			),
			wantErr: false,
		},
		{
			name:     "returns error for invalid JSON",
			jsonData: `{"id":"test-id","invalid-json}`,
			expected: message.Envelope{},
			wantErr:  true,
		},
		{
			name:     "returns error for invalid timestamp format",
			jsonData: `{"id":"test-id","schema":"schema","source":"source","metadata":{},"payload":{"name":"test","value":1},"timestamp":"invalid-timestamp"}`,
			expected: message.Envelope{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var envelope message.Envelope
			err := envelope.UnmarshalJSON([]byte(tt.jsonData))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Payload(), envelope.Payload())
			assert.Equal(t, tt.expected.Base.ID(), envelope.Base.ID())
			assert.Equal(t, tt.expected.Base.Schema(), envelope.Base.Schema())
			assert.Equal(t, tt.expected.Base.Source(), envelope.Base.Source())
			assert.Equal(t, tt.expected.Base.Metadata(), envelope.Base.Metadata())
			assert.True(t, tt.expected.Base.Timestamp().Equal(envelope.Base.Timestamp()))
		})
	}
}

func TestEnvelope_JSONRoundTrip(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	metadata := map[string]string{"key1": "value1", "key2": "value2"}

	original := message.NewEnvelope(
		createTestBaseWithProps(
			"round-trip-id",
			"round-trip-schema",
			"round-trip-source",
			metadata,
			fixedTime,
		),
		map[string]interface{}{"name": "round-trip-test", "value": float64(999)},
	)

	// Marshal to JSON
	jsonData, err := original.MarshalJSON()
	require.NoError(t, err)

	// Unmarshal back to struct
	var unmarshaled message.Envelope
	err = unmarshaled.UnmarshalJSON(jsonData)
	require.NoError(t, err)

	// Verify the round trip preserved all data
	assert.Equal(t, original.Payload(), unmarshaled.Payload())
	assert.Equal(t, original.Base.ID(), unmarshaled.Base.ID())
	assert.Equal(t, original.Base.Schema(), unmarshaled.Base.Schema())
	assert.Equal(t, original.Base.Source(), unmarshaled.Base.Source())
	assert.Equal(t, original.Base.Metadata(), unmarshaled.Base.Metadata())
	assert.True(t, original.Base.Timestamp().Equal(unmarshaled.Base.Timestamp()))
}

func TestEnvelope_WithDifferentPayloadTypes(t *testing.T) {
	t.Run("string payload", func(t *testing.T) {
		envelope := message.NewEnvelope(createTestBase(), "string payload")
		assert.Equal(t, "string payload", envelope.Payload())
	})

	t.Run("int payload", func(t *testing.T) {
		envelope := message.NewEnvelope(createTestBase(), 42)
		assert.Equal(t, 42, envelope.Payload())
	})

	t.Run("map payload", func(t *testing.T) {
		payload := map[string]interface{}{"key": "value", "number": 123}
		envelope := message.NewEnvelope(createTestBase(), payload)
		assert.Equal(t, payload, envelope.Payload())
	})

	t.Run("slice payload", func(t *testing.T) {
		payload := []string{"item1", "item2", "item3"}
		envelope := message.NewEnvelope(createTestBase(), payload)
		assert.Equal(t, payload, envelope.Payload())
	})
}

func TestMarshalEnvelopeJSON(t *testing.T) {
	envelope := message.NewEnvelope(
		createTestBaseWithProps(
			"marshal-id",
			"marshal-schema",
			"marshal-source",
			map[string]string{"meta": "data"},
			time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		),
		TestPayload{Name: "marshal-test", Value: 123},
	)

	jsonData, err := message.MarshalEnvelopeJSON(envelope)
	require.NoError(t, err)

	expectedJSON := `{"id":"marshal-id","schema":"marshal-schema","source":"marshal-source","metadata":{"meta":"data"},"payload":{"name":"marshal-test","value":123},"timestamp":"2024-01-01T12:00:00Z"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func TestUnmarshalEnvelopeJSON(t *testing.T) {
	jsonData := `{"id":"unmarshal-id","schema":"unmarshal-schema","source":"unmarshal-source","metadata":{"meta":"data"},"payload":{"name":"unmarshal-test","value":456},"timestamp":"2024-01-02T15:30:00Z"}`

	envelope, err := message.UnmarshalEnvelopeJSON[TestPayload]([]byte(jsonData))
	require.NoError(t, err)
	payload, ok := message.EnvelopePayloadAs[TestPayload](envelope)
	require.True(t, ok)
	require.NotNil(t, payload)

	// Verify all fields
	assert.Equal(t, "unmarshal-id", envelope.Base.ID())
	assert.Equal(t, "unmarshal-schema", envelope.Base.Schema())
	assert.Equal(t, "unmarshal-source", envelope.Base.Source())
	assert.Equal(t, map[string]string{"meta": "data"}, envelope.Base.Metadata())
	assert.True(t, envelope.Base.Timestamp().Equal(time.Date(2024, 1, 2, 15, 30, 0, 0, time.UTC)))

	assert.Equal(t, "unmarshal-test", payload.Name)
	assert.Equal(t, 456, payload.Value)
}
