package message_test

import (
	"time"
)

type messageMock struct {
}

func (m *messageMock) CausationID() string {
	return "test-causation-id"
}

func (m *messageMock) CorrelationID() string {
	return "test-correlation-id"
}

func (m *messageMock) Metadata() map[string]string {
	return map[string]string{
		"test-key": "test-value",
	}
}

func (m *messageMock) Timestamp() time.Time {
	return time.Now()
}

func (m *messageMock) MessageID() string {
	return "test-message-id"
}
