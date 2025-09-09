package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/message"
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

// This is the event contract for this example.
type SpeechProcessedEvent interface {
	message.Message
	Speech() string
	IsError() bool
}

// A sample event implementation for this example.
type speechProcessedEvent struct {
	event.Base

	speech  string
	isError bool
}

func (e speechProcessedEvent) Speech() string {
	return e.speech
}

func (e speechProcessedEvent) IsError() bool {
	return e.isError
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	// Create an in-memory event bus.
	bus := event.NewInMemoryBus()
	topic := "com.org.speech_event"
	handler := func(ctx context.Context, evt SpeechProcessedEvent) error {
		fmt.Printf("\n📨 %s[Event Handled]: %s %s\n", Green, evt.Speech(), Reset)
		if !evt.IsError() {
			return nil
		}
		return errors.New("❌ Simulating an error handling the event")
	}

	// Register the event handler.
	err := event.Handle(bus, topic, handler)
	if err != nil {
		panicErr(err)
	}

	// Publish a couple of events.
	// The second one will simulate an error.
	// Note that the event handler is executed in a separate goroutine,
	// so the main goroutine can continue to publish events or do other work.
	// In a real application, you might want to wait for all events to be processed
	// before exiting the application.
	if err = publishEvent(ctx, bus, speechProcessedEvent{
		Base:    event.New(),
		speech:  "This is a sample event!",
		isError: false,
	}); err != nil {
		panicErr(err)
	}

	if err = publishEvent(ctx, bus, speechProcessedEvent{
		Base:    event.New(),
		speech:  "This is a sample event simulating an error!",
		isError: true, // Just to simulate an error.
	}); err != nil {
		fmt.Printf("\n%s%s%s\n", Red, err.Error(), Reset)
	}

	cancelSignal()
}

func publishEvent(ctx context.Context, bus event.Bus, evt SpeechProcessedEvent) error {
	fmt.Printf("\n🚀 %s[Publishing Event]: %s%s\n", Cyan, evt.Speech(), Reset)
	return bus.Publish(ctx, "com.org.speech_event", evt)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
