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

// A sample command payload for this example.
type PrintSpeechCommand struct {
	messaging.BaseCommand

	Speech  string `json:"speech"`
	IsError bool   `json:"-"`
}

// PrintSpeechCommandHandler is a sample command handler that processes SpeechCommand.
type PrintSpeechCommandHandler struct {
	wg *sync.WaitGroup
}

func (h PrintSpeechCommandHandler) Handle(_ context.Context, cmd PrintSpeechCommand) error {
	defer h.wg.Done()

	// Simulate an error if IsError is true.
	if cmd.IsError {
		return errors.New("simulated error handling command")
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
	wg.Add(2) // We plan to dispatch 2 commands.

	// Create an in-memory command bus and subscribe a handler to it.
	bus := messaging.NewInMemoryCommandBus()
	unsub, err := messaging.SubscribeCommand(
		rootCtx,
		bus,
		"com.org.test_command",
		PrintSpeechCommandHandler{wg: wg},
	)
	if err != nil {
		panicErr(err)
	}

	// Dispatch a couple of commands.
	if err = dispatchCommand(ctx, bus, PrintSpeechCommand{
		BaseCommand: messaging.NewBaseCommand("com.org.test_command"),
		Speech:      "Welcome to the Command Bus example!",
		IsError:     false,
	}); err != nil {
		panicErr(err)
	}

	if err = dispatchCommand(ctx, bus, PrintSpeechCommand{
		BaseCommand: messaging.NewBaseCommand("com.org.test_command"),
		Speech:      "Let's simulate an error!",
		IsError:     true, // Just to simulate an error.
	}); err != nil {
		//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
		fmt.Printf("❌ %sError dispatching command: %s%s\n", Red, err.Error(), Reset)
	}

	// Unsubscribe the handler and shutdown the bus.
	unsub()

	// Wait for all commands to be processed.
	wg.Wait()
}

func dispatchCommand(ctx context.Context, bus messaging.CommandDispatcher, cmd PrintSpeechCommand) error {
	//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
	fmt.Printf("\n🚀 %s[Dispatching Command]: %s%s\n", Cyan, cmd.Speech, Reset)
	return bus.Dispatch(ctx, cmd)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
