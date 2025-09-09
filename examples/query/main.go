package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/message/query"
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

// A sample query payload for this example.
type SpeechQuery struct {
	query.Base
	Speech   string `json:"speech"`
	IsError  bool   `json:"-"`
	Response any    `json:"response,omitempty"`
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	bus := query.NewInMemoryBus()
	err := query.Handle(bus, "com.org.test_query", func(ctx context.Context, qry SpeechQuery) (any, error) {
		fmt.Printf("\n📨 %s[Query Handled]: %s %s\n", Green, qry.Speech, Reset)
		if qry.IsError {
			return nil, errors.New("❌ Simulating an error processing the query")
		}
		if qry.Response != nil {
			return qry.Response, nil
		}
		return nil, nil
	})
	if err != nil {
		panicErr(err)
	}

	if _, err = dispatchQuery(ctx, bus, SpeechQuery{
		Speech:  "Welcome to the Query Bus example!",
		IsError: false,
	}); err != nil {
		panicErr(err)
	}

	if _, err = dispatchQuery(ctx, bus, SpeechQuery{
		Speech:  "Let's simulate an error!",
		IsError: true, // Just to simulate an error.
	}); err != nil {
		fmt.Printf("\n%s%s%s\n", Red, err.Error(), Reset)
	}

	// Sample query with response
	if res, err := dispatchQuery(ctx, bus, SpeechQuery{
		Speech:   "How are you?",
		Response: "I'm just a computer program, but thanks for asking!",
		IsError:  false,
	}); err != nil {
		panicErr(err)
	} else {
		fmt.Printf("\n%s📩 [Query Response]: %v%s\n", Green, res, Reset)
	}

	cancelSignal()
}

func dispatchQuery(ctx context.Context, bus query.Bus, qry SpeechQuery) (any, error) {
	fmt.Printf("\n-----------\n")
	fmt.Printf("\n🚀 %s[Dispatching Query]: %s%s\n", Cyan, qry.Speech, Reset)
	return bus.Dispatch(ctx, "com.org.test_query", qry)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
