package event_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/xfrr/go-cqrsify/event"
)

type MockEventPayload struct {
	Greeting string
}

type MockResponse struct {
	Result string
}

func MockHandler(ctx event.Context[string, MockEventPayload]) error {
	return nil
}

var (
	bufferSizes = []uint{1, 10, 100, 500, 1000}
)

func BenchmarkBus_Publish(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, bufferSize := range bufferSizes {
		name := fmt.Sprintf("buffer-size-%d", bufferSize)
		b.Run(name, func(b *testing.B) {

			bus, err := event.NewBus(event.WithBufferSize(bufferSize))
			if err != nil {
				panic(err)
			}

			handler := event.NewHandler[string, MockEventPayload](bus)
			_, err = handler.Handle(ctx, "event-reason", MockHandler)
			if err != nil {
				panic(err)
			}

			payload := MockEventPayload{
				Greeting: "Hello World!",
			}

			evt := event.New(
				"event-id",
				"event-reason",
				payload,
			).Any()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := bus.Publish(ctx, evt)
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
