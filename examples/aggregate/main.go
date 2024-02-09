package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/xfrr/cqrsify/aggregate"
	"github.com/xfrr/cqrsify/event"
)

const (
	CustomAggregateCreatedEventName       = "aggregate.created"
	CustomAggregateStatusChangedEventName = "aggregate.status_changed"
)

func makeID() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", rnd.Intn(1000))
}

func main() {
	// create a new aggregate with a random ID
	agg := makeAggregate(makeID(), "aggregate-name")

	fmt.Printf("Aggregate is %s at version %d\n", agg.Status, agg.AggregateVersion())

	// change the status
	agg.ChangeStatus("ready")

	// commit the changes
	agg.CommitChanges()

	fmt.Printf("Aggregate is now %s at version %d\n", agg.Status, agg.AggregateVersion())
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
	*aggregate.Base

	CustomAggregateRoot
}

// ChangeStatus changes the status of the aggregate
// and records the change as an event
func (a *CustomAggregate) ChangeStatus(status string) error {
	// business logic and validation goes here
	// ...

	// use aggregate.Change to apply the event and record it
	aggregate.ApplyChange(a, makeID(), CustomAggregateStatusChangedEventName, CustomAggregateStatusChanged{
		Previous: a.Status,
		New:      status,
	})

	return nil
}

func (a *CustomAggregate) handleStatusChangedEvent(evt event.Event[any]) error {
	payload, ok := evt.Payload().(CustomAggregateStatusChanged)
	if !ok {
		return errors.New("invalid event payload")
	}

	a.changeStatus(payload.New)
	return nil
}

func (a *CustomAggregate) changeStatus(status string) {
	a.Status = status
}

func makeAggregate(id string, name string) *CustomAggregate {
	// create a new aggregate with embedded aggregate.Base
	agg := &CustomAggregate{
		aggregate.New(id, name),
		CustomAggregateRoot{
			Status: "created",
		},
	}

	// apply the change to the aggregate
	aggregate.ApplyChange(
		agg,
		makeID(),
		CustomAggregateCreatedEventName,
		CustomAggregateCreatedEvent{
			ID:     agg.AggregateID().String(),
			Status: agg.Status,
		})

	// handle events
	agg.When(CustomAggregateStatusChangedEventName, agg.handleStatusChangedEvent)
	return agg
}
