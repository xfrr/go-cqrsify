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

type CreateOrderCommand struct {
	messaging.BaseCommand

	OrderID int
}

type CreateOrderCommandJSONPayload struct {
	OrderID int `json:"orderId"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	command := CreateOrderCommand{
		BaseCommand: messaging.NewBaseCommand("com.cqrsify.examples.commands.order.create.v1"),
		OrderID:     123,
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
	unsub, err := messaging.SubscribeCommand(ctx, commandBus, messaging.CommandHandlerFn[CreateOrderCommand](func(ctx context.Context, command CreateOrderCommand) error {
		fmt.Println("Handling command:")
		fmt.Printf("- Command Type: %s\n", command.MessageType())
		fmt.Printf("- Order ID: %d\n", command.OrderID)
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

	cleanup := func() {
		nc.Drain()
		nc.Close()
	}

	js, err := jetstream.New(nc)
	if err != nil {
		cleanup()
		panic(err)
	}

	// Ensure the stream exists
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"com.cqrsify.examples.commands.order.create.v1"},
		Storage:   jetstream.MemoryStorage,
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		cleanup()
		panic(err)
	}

	publisher, err := messagingnats.NewJetStreamMessagePublisher(
		nc,
		streamName,
		serializer,
		deserializer,
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	consumer, err := messagingnats.NewJetStreamMessageConsumer(
		nc,
		streamName,
		serializer,
		deserializer,
		messagingnats.WithConsumerConfig(jetstream.ConsumerConfig{
			Name:          "cqrsify_examples_command_bus_consumer",
			AckPolicy:     jetstream.AckExplicitPolicy,
			FilterSubject: "com.cqrsify.examples.commands.>",
		}),
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	commandBus := messagingnats.NewJetStreamCommandBus(publisher, consumer)
	return commandBus, cleanup, nil
}

func registerCreateOrderCommandJSONSerializer(
	serializer *messaging.JSONSerializer,
	msgType string,
) *messaging.JSONSerializer {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e CreateOrderCommand) messaging.JSONMessage[CreateOrderCommandJSONPayload] {
			return messaging.NewJSONMessage(e, CreateOrderCommandJSONPayload{
				OrderID: e.OrderID,
			})
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
		func(jsonMessage messaging.JSONMessage[CreateOrderCommandJSONPayload]) (CreateOrderCommand, error) {
			return CreateOrderCommand{
				BaseCommand: messaging.NewCommandFromJSON(jsonMessage),
				OrderID:     jsonMessage.Payload.OrderID,
			}, nil
		})
}
