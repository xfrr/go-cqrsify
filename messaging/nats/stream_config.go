package messagingnats

import (
	"github.com/nats-io/nats.go/jetstream"
)

func defaultStreamConfig(streamName string, subjects ...string) jetstream.StreamConfig {
	if len(subjects) == 0 {
		subjects = []string{"cqrsify.messages.>"}
	}

	return jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
	}
}
