package event_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate/event"
)

type MockEventPayload struct {
	Greeting string
}

type MockResponse struct {
	Result string
}

func MockHandler(_ event.Context[string, MockEventPayload]) error {
	return nil
}

func BenchmarkBus_Publish(b *testing.B) {
	var (
		bufferSizes = []uint{1, 10, 100, 500, 1000}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, bufferSize := range bufferSizes {
		name := fmt.Sprintf("buffer-size-%d", bufferSize)
		b.Run(name, func(b *testing.B) {

			bus, err := event.NewInMemoryBus(event.WithBufferSize(bufferSize))
			if err != nil {
				panic(err)
			}

			handler := event.NewHandler[string, MockEventPayload](bus)
			_, err = handler.Handle(ctx, "event-name", MockHandler)
			if err != nil {
				panic(err)
			}

			payload := MockEventPayload{
				Greeting: "Hello World!",
			}

			evt, err := event.New(
				"event-id",
				"event-name",
				payload,
			)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err = bus.Publish(ctx, evt.Any())
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
