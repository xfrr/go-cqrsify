package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"

	messagingnats "github.com/xfrr/go-cqrsify/messaging/nats"
)

const streamName = "cqrsify_command_bus_example"

type CreateOrderCommand interface {
	messaging.Command

	OrderID() int
}

type createOrderCommand struct {
	messaging.BaseCommand

	orderID int
}

func (e createOrderCommand) OrderID() int {
	return e.orderID
}

type createOrderCommandPayload struct {
	OrderID int `json:"orderId"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	command := createOrderCommand{
		BaseCommand: messaging.NewBaseCommand("com.cqrsify.commands.order.create.v1"),
		orderID:     123,
	}

	// register serializers and deserializers
	serializer := messaging.NewJSONSerializer()
	registerCreateOrderCommandJSONSerializer(serializer, command.MessageType())
	deserializer := messaging.NewJSONDeserializer()
	registerCreateOrderCommandJSONDeserializer(deserializer, command.MessageType())

	commandBus, closeCommandBus, err := newCommandBus(
		ctx,
		streamName,
		serializer,
		deserializer,
	)
	if err != nil {
		panic(err)
	}
	defer closeCommandBus()

	// Subscribe to messages of type "CreateOrder"
	unsub, err := messaging.SubscribeCommand(ctx, commandBus, command.MessageType(), messaging.CommandHandlerFn[CreateOrderCommand](func(ctx context.Context, command CreateOrderCommand) error {
		fmt.Println("Handling command:")
		fmt.Printf("- Command Type: %s\n", command.MessageType())
		fmt.Printf("- Order ID: %d\n", command.OrderID())
		return nil
	}))
	if err != nil {
		panic(err)
	}
	defer unsub()

	fmt.Printf("Publishing command '%s'\n", command.MessageType())

	// Publish the command
	err = messaging.DispatchCommand(ctx, commandBus, command)
	if err != nil {
		panic(err)
	}

	// Wait for a while to ensure the command is received before exiting
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
	}
}

func newCommandBus(
	ctx context.Context,
	streamName string,
	serializer *messaging.JSONSerializer,
	deserializer *messaging.JSONDeserializer,
) (messaging.CommandBus, func(), error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	commandBus, err := messagingnats.NewJetStreamCommandBus(
		ctx,
		nc,
		streamName,
		serializer,
		deserializer,
		messagingnats.WithStreamConfig(
			jetstream.StreamConfig{
				Name:      streamName,
				Subjects:  []string{"com.cqrsify.commands.>"},
				MaxAge:    10 * time.Minute,
				Storage:   jetstream.MemoryStorage,
				Retention: jetstream.WorkQueuePolicy,
			},
		),
	)
	if err != nil {
		nc.Close()
		return nil, nil, err
	}

	cleanup := func() {
		nc.Close()
	}
	return commandBus, cleanup, nil
}

func registerCreateOrderCommandJSONSerializer(
	serializer *messaging.JSONSerializer,
	msgType string,
) *messaging.JSONSerializer {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e CreateOrderCommand) messaging.JSONMessage[createOrderCommandPayload] {
			jsonMessage := messaging.NewJSONMessage(e, createOrderCommandPayload{
				OrderID: e.OrderID(),
			})
			return jsonMessage
		})
	return serializer
}

func registerCreateOrderCommandJSONDeserializer(
	deserializer *messaging.JSONDeserializer,
	msgType string,
) {
	messaging.RegisterJSONMessageDeserializer(
		deserializer,
		msgType,
		func(jsonMessage messaging.JSONMessage[createOrderCommandPayload]) (CreateOrderCommand, error) {
			parsedCommand := messaging.NewCommandFromJSON(jsonMessage)
			return &createOrderCommand{
				BaseCommand: parsedCommand,
				orderID:     jsonMessage.Payload.OrderID,
			}, nil
		})
}
