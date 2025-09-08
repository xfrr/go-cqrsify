package message_test

import (
	"time"
)

type messageMock struct {
}

func (m *messageMock) ID() string {
	return "test-message-id"
}

func (m *messageMock) Schema() string {
	return "test-schema"
}

func (m *messageMock) Source() string {
	return "test-source"
}

func (m *messageMock) Metadata() map[string]string {
	return map[string]string{
		"test-key": "test-value",
	}
}

func (m *messageMock) Timestamp() time.Time {
	return time.Now()
}
