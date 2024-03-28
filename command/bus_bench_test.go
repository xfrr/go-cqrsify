package command_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/xfrr/go-cqrsify/command"
)

type MockCommandPayload struct {
	Greeting string
}

type MockResponse struct {
	Result string
}

func MockHandler(ctx command.Context[MockCommandPayload]) error {
	return nil
}

var (
	bufferSizes = []uint{1, 10, 100, 500, 1000}
)

func BenchmarkBus_Dispatch(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, bufferSize := range bufferSizes {
		name := fmt.Sprintf("buffer-size-%d", bufferSize)
		b.Run(name, func(b *testing.B) {

			bus, err := command.NewBus(command.WithBufferSize(bufferSize))
			if err != nil {
				panic(err)
			}

			handler := command.NewHandler[MockCommandPayload](bus)
			_, err = handler.Handle(ctx, "Command", MockHandler)
			if err != nil {
				panic(err)
			}

			payload := MockCommandPayload{
				Greeting: "Hello World!",
			}

			cmd := command.New[MockCommandPayload](
				"command-id",
				payload,
			).Any()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := bus.Dispatch(ctx, "Command", cmd)
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
