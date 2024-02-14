package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/cqrsify/command"
)

const (
	MockCommandID     = "MockCommandID"
	MockAggregateID   = "MockAggregateID"
	MockAggregateName = "MockAggregate"

	MockCommandSubject = "commands.greeting"
)

var (
	done = make(chan struct{})
)

type Payload struct {
	Greeting string `json:"greeting"`
}

func GreetingCommandHandler(ctx command.Context[Payload]) error {
	fmt.Printf("[aggregate_id]: %s\n", ctx.Command().AggregateID())
	fmt.Printf("[aggregate_name]: %s\n", ctx.Command().AggregateName())
	fmt.Printf("[command_id]: %s\n", ctx.Command().ID())
	fmt.Printf("[command_payload]: %s\n", ctx.Command().Payload().Greeting)
	done <- struct{}{}
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bus, err := command.NewBus()
	if err != nil {
		panic(err)
	}

	if err = handleCommands(ctx, bus); err != nil {
		panic(err)
	}

	if err = dispatchGreetingCommand(ctx, bus, Payload{
		Greeting: "Hello World!",
	}); err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
	case <-done:
	}
}

func dispatchGreetingCommand(ctx context.Context, bus command.Bus, payload Payload) error {
	cmd := command.New[Payload](MockCommandID, payload,
		command.WithAggregate(MockAggregateID, MockAggregateName),
	)

	err := bus.Dispatch(ctx, MockCommandSubject, cmd.Any())
	if err != nil {
		return err
	}

	return nil
}

func handleCommands(ctx context.Context, bus command.Bus) error {
	errs, err := command.Handle(ctx, bus, MockCommandSubject, GreetingCommandHandler)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err, ok := <-errs:
				if !ok {
					return
				}
				// handle error...
				fmt.Printf("[error]: %s\n", err.Error())
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
