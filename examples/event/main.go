package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
)

// ANSI escape codes for colors.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// A sample event payload for this example.
type SpeechPrintedEvent struct {
	messaging.BaseEvent

	Speech  string `json:"speech"`
	IsError bool   `json:"-"`
}

// SpeechPrintedEventHandler is a sample event handler that processes SpeechEvent.
type SpeechPrintedEventHandler struct {
	wg *sync.WaitGroup
}

func (h SpeechPrintedEventHandler) Handle(_ context.Context, cmd SpeechPrintedEvent) error {
	defer h.wg.Done()

	// Simulate an error if IsError is true.
	if cmd.IsError {
		return errors.New("simulated error handling event")
	}

	//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
	fmt.Printf("✅ %s%s%s\n", Green, cmd.Speech, Reset)
	return nil
}

func main() {
	rootCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(rootCtx, os.Interrupt)
	defer cancelSignal()

	wg := &sync.WaitGroup{}
	//nolint:mnd // just for this example.
	wg.Add(2) // We plan to publish 2 events.

	// Create an in-memory event bus and subscribe a handler to it.
	bus := messaging.NewInMemoryEventBus()
	unsub, err := messaging.SubscribeEvent(
		rootCtx,
		bus,
		"com.org.speech_printed.v1",
		SpeechPrintedEventHandler{wg: wg},
	)
	if err != nil {
		panicErr(err)
	}

	// Publish a couple of events.
	if err = publishEvent(ctx, bus, SpeechPrintedEvent{
		BaseEvent: messaging.NewBaseEvent("com.org.speech_printed.v1"),
		Speech:    "Welcome to the Event Bus example!",
		IsError:   false,
	}); err != nil {
		panicErr(err)
	}

	if err = publishEvent(ctx, bus, SpeechPrintedEvent{
		BaseEvent: messaging.NewBaseEvent("com.org.speech_printed.v1"),
		Speech:    "Let's simulate an error handling this event.",
		IsError:   true, // Just to simulate an error.
	}); err != nil {
		//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
		fmt.Printf("❌ %sError publishing event: %s%s\n", Red, err.Error(), Reset)
	}

	// Unsubscribe the handler and shutdown the bus.
	unsub()

	// Wait for all events to be processed.
	wg.Wait()
}

func publishEvent(ctx context.Context, bus messaging.EventPublisher, cmd SpeechPrintedEvent) error {
	//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
	fmt.Printf("\n🚀 %s[Publishing Event]: %s%s\n", Cyan, cmd.Speech, Reset)
	return bus.Publish(ctx, cmd)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
