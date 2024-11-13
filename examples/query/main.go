package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/xfrr/go-cqrsify/cqrs"
)

const TimeoutSeconds = 1

// ANSI escape codes for colors.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// A sample query payload for this example.
type GetProductNameQuery struct {
	ProductID string
}

type GetProductNameQueryResponse struct {
	ProductName string
}

func (c GetProductNameQuery) QueryName() string {
	return "get-product-query"
}

// The handler function we will use to handle the GetProductNameQuery.
func GetProductNameQueryHandler(_ context.Context, gpnq GetProductNameQuery) (*GetProductNameQueryResponse, error) {
	fmt.Printf("\n📨 Handler: %sGetting product name by ID: %s %s\n", Green, gpnq.ProductID, Reset)

	if gpnq.ProductID == "123" {
		return &GetProductNameQueryResponse{
			ProductName: "sample-product-name",
		}, nil
	}

	return nil, errors.New("error getting product name")
}

func main() {
	// Wait for interrupt signal to gracefully shutdown the app.
	// Press Ctrl+C to trigger the interrupt signal.
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	bus := cqrs.NewInMemoryBus()
	bus.Use(cqrs.RecoverMiddleware(func(r interface{}) {
		fmt.Printf("\n%s🚨 Recovered from panic: %v%s\n", Red, r, Reset)
	}))
	defer bus.Close()

	err := cqrs.Handle(ctx, bus, GetProductNameQueryHandler)
	if err != nil {
		panicErr(err)
	}

	res, err := dispatchQuery(ctx, bus, GetProductNameQuery{
		ProductID: "123",
	})
	if err != nil {
		panicErr(err)
	}

	fmt.Printf("\n📦 Main: %sThe product name is: %s %s\n", Yellow, res.ProductName, Reset)

	_, err = dispatchQuery(ctx, bus, GetProductNameQuery{
		ProductID: "456",
	})
	if err != nil {
		fmt.Printf("\n🚨 Main: %s %s %s\n", Red, err.Error(), Reset)
	}

	cancelSignal()
}

func dispatchQuery(ctx context.Context, bus cqrs.Bus, qry GetProductNameQuery) (*GetProductNameQueryResponse, error) {
	fmt.Printf("\n🚀 Main: %sDispatching Query: %s%s\n", Cyan, qry.ProductID, Reset)
	return cqrs.Dispatch[*GetProductNameQueryResponse](ctx, bus, qry)
}

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("❌ %s", err.Error()))
	}
}
