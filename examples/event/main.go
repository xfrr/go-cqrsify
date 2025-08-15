package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/message/event"
)

const TimeoutSeconds = 1

// ANSI escape codes for colors.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// A sample event payload for this example.
type SpeechProcessedEvent struct {
	event.Base
	Speech  string `json:"speech"`
	IsError bool   `json:"-"` // just for the sake of example
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	bus := event.NewInMemoryBus()
	err := event.Handle(bus, func(ctx context.Context, evt SpeechProcessedEvent) error {
		fmt.Printf("\n📨 %s[Event Handled]: %s %s\n", Green, evt.Speech, Reset)
		if !evt.IsError {
			return nil
		}
		return errors.New("❌ Simulating an error handling the event")
	})
	if err != nil {
		panicErr(err)
	}

	if err = publishEvent(ctx, bus, SpeechProcessedEvent{
		Speech:  "This is a sample event!",
		IsError: false,
	}); err != nil {
		panicErr(err)
	}

	if err = publishEvent(ctx, bus, SpeechProcessedEvent{
		Speech:  "This is a sample event simulating an error!",
		IsError: true, // Just to simulate an error.
	}); err != nil {
		fmt.Printf("\n%s%s%s\n", Red, err.Error(), Reset)
	}

	cancelSignal()
}

func publishEvent(ctx context.Context, bus event.Bus, evt SpeechProcessedEvent) error {
	fmt.Printf("\n🚀 %s[Publishing Event]: %s%s\n", Cyan, evt.Speech, Reset)
	return bus.Publish(ctx, evt)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
