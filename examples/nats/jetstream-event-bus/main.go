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

const streamName = "cqrsify_event_bus_example"

type OrderCreatedEvent interface {
	messaging.Event

	OrderID() int
}

type orderCreatedEvent struct {
	messaging.BaseEvent

	orderID int
}

func (e orderCreatedEvent) OrderID() int {
	return e.orderID
}

type orderCreatedEventPayload struct {
	OrderID int `json:"orderId"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	event := orderCreatedEvent{
		BaseEvent: messaging.NewBaseEvent("com.cqrsify.examples.events.order.created.v1"),
		orderID:   123,
	}

	// register serializers and deserializers
	serializer := messaging.NewJSONSerializer()
	registerOrderCreatedEventJSONSerializer(serializer, event.MessageType())
	deserializer := messaging.NewJSONDeserializer()
	registerOrderCreatedEventJSONDeserializer(deserializer, event.MessageType())

	eventBus, closeEventBus, err := newEventBus(
		ctx,
		streamName,
		serializer,
		deserializer,
	)
	if err != nil {
		panic(err)
	}
	defer closeEventBus()

	// Subscribe to messages of type "OrderCreated"
	unsub, err := messaging.SubscribeEvent(ctx, eventBus, messaging.EventHandlerFn[OrderCreatedEvent](func(ctx context.Context, event OrderCreatedEvent) error {
		fmt.Println("Handling event:")
		fmt.Printf("- Event Type: %s\n", event.MessageType())
		fmt.Printf("- Order ID: %d\n", event.OrderID())
		return nil
	}))
	if err != nil {
		panic(err)
	}
	defer unsub()

	fmt.Printf("Publishing event '%s'\n", event.MessageType())

	// Publish the event
	err = eventBus.Publish(ctx, event)
	if err != nil {
		panic(err)
	}

	// Wait for a while to ensure the event is received before exiting
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
	}
}

func newEventBus(
	ctx context.Context,
	streamName string,
	serializer *messaging.JSONSerializer,
	deserializer *messaging.JSONDeserializer,
) (messaging.EventBus, func(), error) {
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
		Subjects:  []string{"com.cqrsify.examples.events.>"},
		Storage:   jetstream.MemoryStorage,
		Retention: jetstream.LimitsPolicy,
	})
	if err != nil {
		cleanup()
		panic(err)
	}

	eventPublisher, err := messagingnats.NewJetStreamMessagePublisher(
		nc,
		streamName,
		serializer,
		deserializer,
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	eventConsumer, err := messagingnats.NewJetStreamMessageConsumer(
		ctx,
		nc,
		streamName,
		serializer,
		deserializer,
		messagingnats.WithJetStreamConsumerConfig(jetstream.ConsumerConfig{
			Name:          "cqrsify_examples_event_bus_consumer",
			AckPolicy:     jetstream.AckExplicitPolicy,
			FilterSubject: "com.cqrsify.examples.events.>",
		}),
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	eventBus := messagingnats.NewJetStreamEventBus(eventPublisher, eventConsumer)
	return eventBus, cleanup, nil
}

func registerOrderCreatedEventJSONSerializer(
	serializer *messaging.JSONSerializer,
	msgType string,
) *messaging.JSONSerializer {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e OrderCreatedEvent) messaging.JSONMessage[orderCreatedEventPayload] {
			jsonMessage := messaging.NewJSONMessage(e, orderCreatedEventPayload{
				OrderID: e.OrderID(),
			})
			return jsonMessage
		})
	return serializer
}

func registerOrderCreatedEventJSONDeserializer(
	deserializer *messaging.JSONDeserializer,
	msgType string,
) {
	messaging.RegisterJSONMessageDeserializer(
		deserializer,
		msgType,
		func(jsonMessage messaging.JSONMessage[orderCreatedEventPayload]) (OrderCreatedEvent, error) {
			parsedEvent := messaging.NewEventFromJSON(jsonMessage)
			return &orderCreatedEvent{
				BaseEvent: parsedEvent,
				orderID:   jsonMessage.Payload.OrderID,
			}, nil
		})
}
