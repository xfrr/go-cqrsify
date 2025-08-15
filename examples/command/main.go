package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/message/command"
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
	command.Base
	Speech  string `json:"speech"`
	IsError bool   `json:"-"`
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	bus := command.NewInMemoryBus()
	err := command.Handle(bus, func(ctx context.Context, cmd SpeechCommand) error {
		fmt.Printf("\n📨 %s[Command Handled]: %s %s\n", Green, cmd.Speech, Reset)
		if !cmd.IsError {
			return nil
		}
		return errors.New("❌ Simulating an error processing the command")
	})
	if err != nil {
		panicErr(err)
	}

	if err = dispatchCommand(ctx, bus, SpeechCommand{
		Speech:  "This is a sample speech command!",
		IsError: false,
	}); err != nil {
		panicErr(err)
	}

	if err = dispatchCommand(ctx, bus, SpeechCommand{
		Speech:  "This is a sample speech command simulating an error!",
		IsError: true, // Just to simulate an error.
	}); err != nil {
		fmt.Printf("\n%s%s%s\n", Red, err.Error(), Reset)
	}

	cancelSignal()
}

func dispatchCommand(ctx context.Context, bus command.Bus, cmd SpeechCommand) error {
	fmt.Printf("\n🚀 %s[Dispatching Command]: %s%s\n", Cyan, cmd.Speech, Reset)
	return bus.Dispatch(ctx, cmd)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
