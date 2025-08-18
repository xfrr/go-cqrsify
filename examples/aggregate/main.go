package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"
)

// sample event names
const (
	CustomAggregateCreatedEventName      = "aggregate.created"
	CustomAggregateStatusEventdEventName = "aggregate.status_event"
)

func main() {
	// create a new aggregate with a random ID
	agg := makeAggregate(randomStr(), "aggregate-name")

	log.Printf("Aggregate created: %s\n", coloured(agg.String()))

	// event the aggregate status and commit the event
	eventStatus(agg, "ready")

	log.Printf("Aggregate status eventd: %s\n", coloured(agg.String()))
}

type CustomAggregateStatusEventd struct {
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

func (agg *CustomAggregate) EventStatus(status string) error {
	// business logic and validation goes here
	// ...

	// generate a new random event ID
	eventID := randomStr()

	// use aggregate.RaiseEvent to apply the event to the aggregate
	err := aggregate.RaiseEvent(
		agg,
		eventID,
		CustomAggregateStatusEventdEventName,
		CustomAggregateStatusEventd{
			Previous: agg.Status,
			New:      status,
		},
	)
	if err != nil {
		panic(err)
	}

	return nil
}

func (agg *CustomAggregate) handleStatusEventdEvent(e event.Event[any, any]) {
	evt, ok := event.Cast[string, CustomAggregateStatusEventd](e)
	if !ok {
		log.Fatalf("failed to cast event %s to CustomAggregateStatusEventd\n", e.Name())
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

	// start listening for status event events
	agg.HandleEvent(CustomAggregateStatusEventdEventName, agg.handleStatusEventdEvent)

	// apply the event to the aggregate
	aggregate.RaiseEvent(
		agg,
		randomStr(),
		CustomAggregateCreatedEventName,
		CustomAggregateCreatedEvent{
			ID:     agg.AggregateID(),
			Status: agg.Status,
		})

	// commit the event
	agg.CommitEvents()
	return agg
}

func eventStatus(agg *CustomAggregate, status string) {
	// apply the event to the aggregate
	agg.EventStatus(status)

	// commit the event
	agg.CommitEvents()
}

func coloured(s string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", s)
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
