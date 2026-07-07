package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/xfrr/go-cqrsify/messaging"
	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

const (
	sampleQueryType      = "cqrsify.queries.find_greeting"
	sampleQueryReplyType = sampleQueryType + ".reply"
)

type sampleQuery struct {
	messaging.BaseQuery

	sampleQueryAttributes
}

type sampleQueryAttributes struct {
	Name string `json:"name"`
}

type sampleQueryReply struct {
	messaging.BaseQueryReply

	Greeting string `json:"greeting"`
}

type sampleQueryHandler struct{}

func (h *sampleQueryHandler) Handle(_ context.Context, qry sampleQuery) (sampleQueryReply, error) {
	return sampleQueryReply{
		BaseQueryReply: messaging.NewMessage(
			sampleQueryReplyType,
			messaging.WithID(qry.QueryID()),
			messaging.WithSource("cqrsify.examples.query-http-server"),
		),
		Greeting: fmt.Sprintf("Hello, %s!", qry.Name),
	}, nil
}

func main() {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	queryBus := messaging.NewInMemoryQueryBus(
		messaging.ConfigureInMemoryMessageBusSubjects(sampleQueryType),
		messaging.ConfigureInMemoryMessageBusErrorHandler(
			func(messageType string, err error) {
				fmt.Printf("Error handling message of type %s: %v\n", messageType, err)
			},
		),
	)

	unsub, err := messaging.RegisterQueryHandler(ctx, queryBus, &sampleQueryHandler{})
	if err != nil {
		panic(err)
	}
	defer unsub()

	queryHandler := messaginghttp.NewQueryHandler(queryBus)
	err = messaginghttp.RegisterJSONAPIQueryDecoder(
		queryHandler,
		sampleQueryType,
		func(_ context.Context, sd apix.SingleDocument[sampleQueryAttributes]) (messaging.Query, error) {
			qry := sampleQuery{
				BaseQuery:             messaginghttp.CreateBaseQueryFromSingleDocument(sampleQueryType, sd),
				sampleQueryAttributes: sd.Data.Attributes,
			}
			return qry, nil
		},
	)
	if err != nil {
		panic(err)
	}

	err = messaginghttp.RegisterSingleDocumentMessageEncoder(
		queryHandler,
		sampleQueryReplyType,
		func(_ context.Context, reply sampleQueryReply) (apix.SingleDocument[sampleQueryReply], error) {
			return apix.NewSingleDocument(reply.MessageType(), reply.MessageID(), reply), nil
		},
	)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	queryHTTPServer := messaginghttp.NewQueryGINServer(queryHandler, r)

	go func() {
		_ = queryHTTPServer.ListenAndServe(":8092")
	}()

	defer func() {
		_ = queryHTTPServer.Close()
	}()

	fmt.Println("HTTP Query Server is running on :8092")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println(" ")
	fmt.Println("Example curl command to send a query with the new QUERY method:")
	fmt.Println(`> curl -X QUERY http://localhost:8092/queries -H "Content-Type: application/vnd.api+json" -d '{"data": {"type": "cqrsify.queries.find_greeting", "id": "qry-123", "attributes": {"name": "CQRS"}}}'`)
	fmt.Println(" ")
	fmt.Println("Equivalent POST request for older clients:")
	fmt.Println(`> curl -X POST http://localhost:8092/queries -H "Content-Type: application/vnd.api+json" -d '{"data": {"type": "cqrsify.queries.find_greeting", "id": "qry-123", "attributes": {"name": "CQRS"}}}'`)

	<-ctx.Done()
	fmt.Println("Shutting down HTTP Query Server...")
}