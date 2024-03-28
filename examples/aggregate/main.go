package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/event"
)

// sample event names
const (
	CustomAggregateCreatedEventName       = "aggregate.created"
	CustomAggregateStatusChangedEventName = "aggregate.status_changed"
)

func main() {
	// create a new aggregate with a random ID
	agg := makeAggregate(randomStr(), "aggregate-name")

	log.Printf("Aggregate created: %s\n", coloured(agg.String()))

	// change the aggregate status and commit the change
	changeStatus(agg, "ready")

	log.Printf("Aggregate status changed: %s\n", coloured(agg.String()))
}

type CustomAggregateStatusChanged struct {
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
	// embed the aggregate.Base to provide the basic functionality of an aggregate
	*aggregate.Base[string]

	CustomAggregateRoot
}

func (agg *CustomAggregate) String() string {
	return fmt.Sprintf("{ID: %s, Status: %s, Version: %d}", agg.AggregateID(), agg.Status, agg.AggregateVersion())
}

func (agg *CustomAggregate) ChangeStatus(status string) error {
	// business logic and validation goes here
	// ...

	// generate a new random event ID
	eventID := randomStr()

	// use aggregate.Change to apply the event and record it
	aggregate.NextChange(
		agg,
		eventID,
		CustomAggregateStatusChangedEventName,
		CustomAggregateStatusChanged{
			Previous: agg.Status,
			New:      status,
		},
	)

	return nil
}

func (agg *CustomAggregate) handleStatusChangedEvent(e event.Event[any, any]) {
	evt, ok := event.Cast[string, CustomAggregateStatusChanged](e)
	if !ok {
		log.Fatalf("failed to cast event %s to CustomAggregateStatusChanged\n", e.Reason())
	}

	agg.Status = evt.Payload().New
}

func makeAggregate(id string, name string) *CustomAggregate {
	// create a new aggregate with embedded aggregate.Base
	agg := &CustomAggregate{
		aggregate.New(id, name),
		CustomAggregateRoot{
			Status: "created",
		},
	}

	log.Printf("Aggregate initialized: %s\n", coloured(agg.String()))

	// start listening for status change events
	agg.When(CustomAggregateStatusChangedEventName, agg.handleStatusChangedEvent)

	// apply the change to the aggregate
	aggregate.NextChange(
		agg,
		randomStr(),
		CustomAggregateCreatedEventName,
		CustomAggregateCreatedEvent{
			ID:     agg.AggregateID(),
			Status: agg.Status,
		})

	// commit the change
	agg.CommitChanges()
	return agg
}

func changeStatus(agg *CustomAggregate, status string) {
	// apply the change to the aggregate
	agg.ChangeStatus(status)

	// commit the change
	agg.CommitChanges()
}

func coloured(s string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", s)
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
