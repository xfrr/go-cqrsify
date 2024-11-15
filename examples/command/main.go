package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/cqrs"
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

// A sample command payload for this example.
type SpeechCommand struct {
	Speech  string `json:"speech"`
	IsError bool   `json:"-"`
}

func (c SpeechCommand) CommandName() string {
	return "speech-command"
}

// The handler function we will use to handle the SpeechCommand.
func SpeechCommandHandler(ctx context.Context, cmd SpeechCommand) (cqrs.EmptyRequestResponse, error) {
	fmt.Printf("\n📨 %sCommand Received!: %s %s\n", Green, cmd.Speech, Reset)

	if !cmd.IsError {
		return cqrs.EmptyRequestResponse{}, nil
	}

	return cqrs.EmptyRequestResponse{}, errors.New("simulating an error processing the command")
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	bus := cqrs.NewInMemoryBus()
	bus.Use(cqrs.RecoverMiddleware(func(r interface{}) {
		fmt.Printf("\n%s🚨 Recovered from panic: %v%s\n", Red, r, Reset)
	}))
	defer bus.Close()

	err := cqrs.Handle(ctx, bus, SpeechCommandHandler)
	if err != nil {
		panicErr(err)
	}

	if _, err = dispatchCommand(ctx, bus, SpeechCommand{
		Speech:  "This is a sample speech command!",
		IsError: false,
	}); err != nil {
		panicErr(err)
	}

	if _, err = dispatchCommand(ctx, bus, SpeechCommand{
		Speech:  "This is a sample speech command simulating an error!",
		IsError: true, // Just to simulate an error.
	}); err != nil {
		fmt.Printf("\n%s%s%s\n", Red, err.Error(), Reset)
	}

	cancelSignal()
}

func dispatchCommand(ctx context.Context, bus cqrs.Bus, cmd SpeechCommand) (any, error) {
	fmt.Printf("\n🚀 %sDispatching Command: %s%s\n", Cyan, cmd.Speech, Reset)
	return cqrs.Dispatch[cqrs.EmptyRequestResponse](ctx, bus, cmd)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
