package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xfrr/go-cqrsify/domain"
)

var (
	// ANSI escape codes for colors.
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// This is the domain event contract for this example.
type SpeechProcessedDomainEvent interface {
	domain.Event
	Speech() string
	IsError() bool
}

// A sample domain event implementation for this example.
type speechProcessedDomainEvent struct {
	domain.BaseEvent

	speech  string
	isError bool // to simulate an error in the handler
}

func (e speechProcessedDomainEvent) Speech() string {
	return e.speech
}

func (e speechProcessedDomainEvent) IsError() bool {
	return e.isError
}

type SpeechProcessedEventHandler struct{}

type CustomAggregateRoot struct {
	SpeechText string
}

type CustomAggregate struct {
	// embed the domain.Base to provide the basic functionality of an aggregate
	*domain.BaseAggregate[string]

	CustomAggregateRoot
}

func (h SpeechProcessedEventHandler) Handle(ctx context.Context, evt SpeechProcessedDomainEvent) error {
	fmt.Printf("%s🗣️  [User Registered Event]: %s%s\n", Green, evt.Speech(), Reset)
	if evt.IsError() {
		return fmt.Errorf("simulated handler error for speech: %s", evt.Speech())
	}
	return nil
}

func main() {
	// Create an example aggregate
	agg := &CustomAggregate{
		BaseAggregate: domain.NewAggregate("agg-1", "agg-name"),
		CustomAggregateRoot: CustomAggregateRoot{
			SpeechText: "Hello, World!",
		},
	}

	// Choose sync or async (here: async with 4 workers)
	bus := domain.NewInMemoryEventBus(
		domain.ConfigureEventBusAsyncWorkers(4),
		domain.ConfigureEventBusQueue(1024),
		domain.ConfigureEventBusErrorHandler(func(evt string, err error) {
			// Handle errors here.
			fmt.Printf("\n%s❌ [Error Handling Event]: %s - %s%s\n", Red, evt, err.Error(), Reset)
		}),
	)
	defer bus.Close()

	// Add middleware to the bus
	bus.Use(
		domain.RecoverMiddleware(func(r any) {
			fmt.Printf("\n%s💥 [Panic Recovered]: %v%s\n", Red, r, Reset)
		}),
		domain.TimeoutMiddleware(5000*time.Millisecond),
		domain.RetryBackoffMiddleware(3, 100*time.Millisecond),
	)

	// Register event handlers
	unsub := domain.SubscribeEvent(bus, "com.org.speech_processed", SpeechProcessedEventHandler{})
	defer unsub()

	// Publish many events
	for i := range 10 {
		event := speechProcessedDomainEvent{
			BaseEvent: domain.NewEvent("com.org.speech_processed", domain.CreateEventAggregateRef(agg)),
			speech:    fmt.Sprintf("%s - message %d", agg.SpeechText, i+1),
		}
		if err := bus.Publish(context.Background(), event); err != nil {
			fmt.Printf("%s❌ [Error Publishing Event]: %s%s\n", Red, err.Error(), Reset)
		}
	}

	// Wait a moment to allow async processing (not needed for sync bus)
	time.Sleep(2 * time.Second)
}
