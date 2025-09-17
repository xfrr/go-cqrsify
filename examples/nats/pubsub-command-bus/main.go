package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"

	messagingnats "github.com/xfrr/go-cqrsify/messaging/nats"
)

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
	OrderID int `json:"order_id"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	command := createOrderCommand{
		BaseCommand: messaging.NewBaseCommand("com.example.order.create.v1"),
		orderID:     123,
	}

	// register serializers and deserializers
	serializer := messaging.NewJSONSerializer()
	registerCreateOrderCommandJsonSerializer(serializer, command.MessageType())
	deserializer := messaging.NewJSONDeserializer()
	registerCreateOrderCommandJsonDeserializer(deserializer, command.MessageType())
	commandBus, closeCommandBus, err := newCommandBus(serializer, deserializer)
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
	// Publish a command of type "CreateOrder"
	err = messaging.DispatchCommand[CreateOrderCommand](ctx, commandBus, command)
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
	serializer *messaging.JSONSerializer,
	deserializer *messaging.JSONDeserializer,
) (messaging.CommandBus, func(), error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	commandBus := messagingnats.NewPubSubCommandBus(nc, serializer, deserializer)
	cleanup := func() {
		nc.Close()
	}
	return commandBus, cleanup, nil
}

func registerCreateOrderCommandJsonSerializer(serializer *messaging.JSONSerializer, msgType string) *messaging.JSONSerializer {
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

func registerCreateOrderCommandJsonDeserializer(deserializer *messaging.JSONDeserializer, msgType string) {
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
