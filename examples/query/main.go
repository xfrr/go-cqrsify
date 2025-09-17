package main

import (
	"context"
	"fmt"
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

type GetSpeechQuery interface {
	messaging.Query
}

// A sample query payload for this example.
type getSpeechQuery struct {
	messaging.BaseQuery
}

type getSpeechQueryReply struct {
	messaging.BaseQuery

	Speech string
}

// GetSpeechQueryHandler is a sample query handler that processes SpeechQuery.
type GetSpeechQueryHandler struct {
	wg *sync.WaitGroup
}

func (h GetSpeechQueryHandler) Handle(_ context.Context, query GetSpeechQuery) error {
	defer h.wg.Done()

	// Reply to the query to acknowledge successful handling.
	if err := query.Reply(context.Background(), getSpeechQueryReply{
		BaseQuery: messaging.NewBaseQuery("com.org.test_query.reply"),
		Speech:    "Welcome to the Query Bus example!",
	}); err != nil {
		return fmt.Errorf("failed to reply to query: %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1) // We plan to dispatch 1 query.

	// Create an in-memory query bus and subscribe a handler to it.
	bus := messaging.NewInMemoryQueryBus()
	unsub, err := messaging.SubscribeQuery(
		ctx,
		bus,
		"com.org.test_query",
		GetSpeechQueryHandler{wg: wg},
	)
	if err != nil {
		panicErr(err)
	}

	// Dispatch a couple of querys.
	res, err := dispatchQuery(ctx, bus, getSpeechQuery{
		BaseQuery: messaging.NewBaseQuery("com.org.test_query"),
	})
	if err != nil {
		panicErr(err)
	}

	//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
	fmt.Printf("🎉 %sQuery handled successfully, reply: %+v%s\n", Yellow, res.Speech, Reset)

	// Unsubscribe the handler and shutdown the bus.
	unsub()

	wg.Wait() // Wait for all handlers to complete.
}

func dispatchQuery(ctx context.Context, bus messaging.QueryDispatcher, qry GetSpeechQuery) (getSpeechQueryReply, error) {
	//nolint:forbidigo // Using fmt.Printf for simplicity in this example.
	fmt.Printf("\n🚀 %sDispatching GetSpeech Query%s\n", Cyan, Reset)
	res, err := messaging.DispatchQuery[GetSpeechQuery, getSpeechQueryReply](ctx, bus, qry)
	if err != nil {
		return getSpeechQueryReply{}, err
	}

	return res, nil
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
