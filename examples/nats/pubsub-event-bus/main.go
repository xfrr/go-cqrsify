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

type OrderCreatedEvent interface {
	messaging.Event

	OrderID() int
	OrderAmount() float64
}

type orderCreatedEvent struct {
	messaging.BaseEvent

	orderID     int
	orderAmount float64
}

func (e *orderCreatedEvent) OrderID() int {
	return e.orderID
}

func (e *orderCreatedEvent) OrderAmount() float64 {
	return e.orderAmount
}

type OrderCreatedEventPayload struct {
	OrderID     int     `json:"order_id"`
	OrderAmount float64 `json:"order_amount"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	evt := &orderCreatedEvent{
		BaseEvent:   messaging.NewBaseEvent("com.example.order.created.v1"),
		orderID:     123,
		orderAmount: 456.78,
	}

	serializer := messaging.NewJSONSerializer()
	registerOrderCreatedEventJsonSerializer(serializer, evt.MessageType())

	deserializer := messaging.NewJSONDeserializer()
	registerOrderCreatedEventJsonDeserializer(deserializer, evt.MessageType())

	// Create a NATS-based PubSubMessagePublisher
	publisher, err := messagingnats.NewPubSubMessagePublisher(
		nc,
		serializer,
		deserializer,
	)
	if err != nil {
		nc.Close()
		panic(err)
	}

	// Create a NATS-based PubSubMessageConsumer
	consumer, err := messagingnats.NewPubSubMessageConsumer(
		nc,
		serializer,
		deserializer,
		messagingnats.WithPubSubConsumerSubject("com.example.order.created.v1"),
	)
	if err != nil {
		nc.Close()
		panic(err)
	}

	// Create a NATS-based PubSubMessageBus
	pubSubBus := messagingnats.NewPubSubEventBus(publisher, consumer)

	// Subscribe to messages of type "OrderCreated"
	unsub, err := messaging.SubscribeEvent(
		ctx,
		pubSubBus,
		messaging.MessageHandlerFn[OrderCreatedEvent](func(ctx context.Context, evt OrderCreatedEvent) error {
			fmt.Println("Received event:")
			fmt.Printf("- Event Type: %s\n", evt.MessageType())
			fmt.Printf("- Order ID: %d\n", evt.OrderID())
			fmt.Printf("- Order Amount: %.2f\n", evt.OrderAmount())
			return nil
		}))
	if err != nil {
		panic(err)
	}
	defer unsub()

	fmt.Printf("Publishing event '%s'\n", evt.MessageType())
	// Publish a event of type "OrderCreated"
	if err := pubSubBus.Publish(ctx, evt); err != nil {
		panic(err)
	}

	// Wait for a while to ensure the event is received before exiting
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
	}
}

func registerOrderCreatedEventJsonSerializer(serializer *messaging.JSONSerializer, msgType string) *messaging.JSONSerializer {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e OrderCreatedEvent) messaging.JSONMessage[OrderCreatedEventPayload] {
			jsonMessage := messaging.NewJSONMessage(e, OrderCreatedEventPayload{
				OrderID: e.OrderID(),
			})
			return jsonMessage
		})
	return serializer
}

func registerOrderCreatedEventJsonDeserializer(deserializer *messaging.JSONDeserializer, msgType string) {
	messaging.RegisterJSONMessageDeserializer(
		deserializer,
		msgType,
		func(jsonMessage messaging.JSONMessage[OrderCreatedEventPayload]) (OrderCreatedEvent, error) {
			parsedEvent := messaging.NewEventFromJSON(jsonMessage)
			return &orderCreatedEvent{
				BaseEvent:   parsedEvent,
				orderID:     jsonMessage.Payload.OrderID,
				orderAmount: jsonMessage.Payload.OrderAmount,
			}, nil
		})
}
