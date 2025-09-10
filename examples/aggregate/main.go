package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/xfrr/go-cqrsify/domain"
)

// sample event names
const (
	CustomAggregateStatusEventName = "status_changed"
)

func main() {
	// create a new aggregate with a random ID
	customAggregate := makeAggregate(randomStr(), "aggregate-name")

	log.Printf("Aggregate initialized: %s\n", coloured(customAggregate.String()))

	eventStatus(customAggregate, "created")

	log.Printf("Aggregate created: %s\n", coloured(customAggregate.String()))

	// event the aggregate status and commit the event
	eventStatus(customAggregate, "ready")

	log.Printf("Aggregate ready: %s\n", coloured(customAggregate.String()))
}

type CustomAggregateStatusChangedEvent struct {
	domain.BaseEvent

	Previous string
	New      string
}

type CustomAggregateCreatedEvent struct {
	ID     string
	Status string
}

type CustomAggregateRoot struct {
	Status string
}

type CustomAggregate struct {
	// embed the domain.Base to provide the basic functionality of an aggregate
	*domain.BaseAggregate[string]

	CustomAggregateRoot
}

func (agg *CustomAggregate) String() string {
	return fmt.Sprintf("{ID: %s, Status: %s, Version: %d}", agg.AggregateID(), agg.Status, agg.AggregateVersion())
}

func (agg *CustomAggregate) EventStatus(status string) error {
	// business logic and validation goes here
	// ...

	// apply the event to the aggregate
	err := domain.NextEvent(
		agg,
		CustomAggregateStatusChangedEvent{
			BaseEvent: domain.NewEvent(
				CustomAggregateStatusEventName,
				domain.CreateEventAggregateRef(agg),
			),
			Previous: agg.Status,
			New:      status,
		},
	)
	if err != nil {
		panic(err)
	}

	return nil
}

func handleStatusEvent(agg *CustomAggregate, e CustomAggregateStatusChangedEvent) error {
	agg.Status = e.New
	return nil
}

func makeAggregate(id string, name string) *CustomAggregate {
	// create a new aggregate with embedded domain.Base
	customAggregate := &CustomAggregate{
		domain.NewAggregate(id, name),
		CustomAggregateRoot{
			Status: "init",
		},
	}

	// start listening for status event events
	domain.HandleEvent(
		customAggregate,
		CustomAggregateStatusEventName,
		handleStatusEvent,
	)

	// apply the event to the aggregate
	domain.NextEvent(
		customAggregate,
		CustomAggregateStatusChangedEvent{
			BaseEvent: domain.NewEvent(
				CustomAggregateStatusEventName,
				domain.CreateEventAggregateRef(customAggregate)),
			Previous: "",
			New:      customAggregate.Status,
		})

	// commit the event
	customAggregate.CommitEvents()
	return customAggregate
}

func eventStatus(customAggregate *CustomAggregate, status string) {
	// apply the event to the aggregate
	customAggregate.EventStatus(status)

	// commit the event
	customAggregate.CommitEvents()
}

func coloured(s string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", s)
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
