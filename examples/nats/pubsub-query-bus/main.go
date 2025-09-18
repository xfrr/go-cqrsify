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

type GetOrderAmountQuery interface {
	messaging.Query

	OrderID() int
}

type getOrderAmountQuery struct {
	messaging.BaseQuery

	orderID int
}

func (e getOrderAmountQuery) OrderID() int {
	return e.orderID
}

type getOrderAmountQueryPayload struct {
	OrderID int `json:"order_id"`
}

type getOrderAmountQueryReplyPayload struct {
	OrderAmount float64 `json:"order_amount"`
}

type getOrderAmountQueryReply struct {
	messaging.BaseQueryReply
	getOrderAmountQueryReplyPayload
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	query := getOrderAmountQuery{
		BaseQuery: messaging.NewBaseQuery("com.example.order.get_amount.v1"),
		orderID:   123,
	}

	queryReply := getOrderAmountQueryReply{
		BaseQueryReply: messaging.NewBaseQueryReply(query),
		getOrderAmountQueryReplyPayload: getOrderAmountQueryReplyPayload{
			OrderAmount: 42.50, // Just a dummy amount
		},
	}

	// register serializers and deserializers
	serializer := messaging.NewJSONSerializer()
	registerGetOrderAmountQueryJsonSerializer(serializer, query.MessageType())
	registerGetOrderAmountQueryReplyJsonSerializer(serializer, queryReply.MessageType())
	deserializer := messaging.NewJSONDeserializer()
	registerGetOrderAmountQueryJsonDeserializer(deserializer, query.MessageType())
	registerGetOrderAmountQueryReplyJsonDeserializer(deserializer, queryReply.MessageType())

	queryBus, closeQueryBus, err := newQueryBus(serializer, deserializer)
	if err != nil {
		panic(err)
	}
	defer closeQueryBus()

	// Subscribe to messages of type "GetOrderAmount"
	unsub, err := messaging.SubscribeQuery(ctx, queryBus, query.MessageType(), messaging.QueryHandlerFn[GetOrderAmountQuery](func(ctx context.Context, query GetOrderAmountQuery) error {
		fmt.Println("Handling query:")
		fmt.Printf("- Query Type: %s\n", query.MessageType())
		fmt.Printf("- Order ID: %d\n", query.OrderID())

		// Reply to the query
		if err := query.Reply(ctx, queryReply); err != nil {
			return fmt.Errorf("failed to reply to query: %w", err)
		}

		return nil
	}))
	if err != nil {
		panic(err)
	}
	defer unsub()

	fmt.Printf("Publishing query '%s'\n", query.MessageType())
	// Publish a query of type "GetOrderAmount"
	reply, err := messaging.DispatchQuery[GetOrderAmountQuery, getOrderAmountQueryReply](ctx, queryBus, query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received reply:")
	fmt.Printf("- Reply Type: %s\n", reply.MessageType())
	fmt.Printf("- Order Amount: %.2f\n", reply.OrderAmount)

	// Wait for a while to ensure the query is received before exiting
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
	}
}

func newQueryBus(
	serializer *messaging.JSONSerializer,
	deserializer *messaging.JSONDeserializer,
) (messaging.QueryBus, func(), error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	queryBus := messagingnats.NewPubSubQueryBus(nc, serializer, deserializer)
	cleanup := func() {
		nc.Close()
	}
	return queryBus, cleanup, nil
}

func registerGetOrderAmountQueryJsonSerializer(serializer *messaging.JSONSerializer, msgType string) *messaging.JSONSerializer {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e GetOrderAmountQuery) messaging.JSONMessage[getOrderAmountQueryPayload] {
			jsonMessage := messaging.NewJSONMessage(e, getOrderAmountQueryPayload{
				OrderID: e.OrderID(),
			})
			return jsonMessage
		})
	return serializer
}

func registerGetOrderAmountQueryJsonDeserializer(deserializer *messaging.JSONDeserializer, msgType string) {
	messaging.RegisterJSONMessageDeserializer(
		deserializer,
		msgType,
		func(jsonMessage messaging.JSONMessage[getOrderAmountQueryPayload]) (GetOrderAmountQuery, error) {
			parsedQuery := messaging.NewQueryFromJSON(jsonMessage)
			return &getOrderAmountQuery{
				BaseQuery: parsedQuery,
				orderID:   jsonMessage.Payload.OrderID,
			}, nil
		})
}

func registerGetOrderAmountQueryReplyJsonSerializer(serializer *messaging.JSONSerializer, msgType string) {
	messaging.RegisterJSONMessageSerializer(
		serializer,
		msgType,
		func(e getOrderAmountQueryReply) messaging.JSONMessage[getOrderAmountQueryReplyPayload] {
			jsonMessage := messaging.NewJSONMessage(e, getOrderAmountQueryReplyPayload{
				OrderAmount: e.OrderAmount,
			})
			return jsonMessage
		})
}

func registerGetOrderAmountQueryReplyJsonDeserializer(deserializer *messaging.JSONDeserializer, msgType string) {
	messaging.RegisterJSONMessageDeserializer(
		deserializer,
		msgType,
		func(jsonMessage messaging.JSONMessage[getOrderAmountQueryReplyPayload]) (getOrderAmountQueryReply, error) {
			parsedQueryReply := messaging.NewQueryReplyFromJSON(jsonMessage)
			return getOrderAmountQueryReply{
				BaseQueryReply:                  parsedQueryReply,
				getOrderAmountQueryReplyPayload: jsonMessage.Payload,
			}, nil
		})
}
