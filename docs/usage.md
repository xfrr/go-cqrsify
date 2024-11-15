
# Usage
This section explains how to use go-cqrsify in your application. You can find examples of how to use the library in the [examples](https://github.com/xfrr/go-cqrsify/tree/main/examples).

## Aggregate
An aggregate is a domain-driven design (DDD) concept that represents a cluster of domain objects that can be treated as a single unit. An aggregate will have a unique identifier and a set of events that represent events to the aggregate.

> **Note**: An Event is called a Event in the Aggregate context and is represented by the `Event` struct in the `aggregate` package.

### Creating an aggregate
To create an aggregate, you need to create a struct that embeds the `*aggregate.Base` struct and implements the `aggregate.Aggregate` interface. The `Base` struct provides the basic functionality for an aggregate, such as the `ApplyEvent` method, which is used to apply a event to the aggregate.

```go
type Customer struct {
    *aggregate.Base

    name string
}

func NewCustomer(id string, name string) *Customer {
   // the unique name of the aggregate
   aggregateName := "customer"

    // create a new customer aggregate with embedded aggregate.Base
    customer := &Customer{
        Base: aggregate.New(id, aggregateName),
        name: name,
    }

    // apply the event to the aggregate
    // note: this is a helper function that creates an event and applies it to the 
    // aggregate. You can also create the event manually and apply it to the aggregate 
    // by calling customer.ApplyEvent(event)
	aggregate.RaiseEvent(
		customer,
		CustomAggregateCreatedEventName,
		CustomAggregateCreatedEvent{
			ID:   customer.AggregateID().String(),
			Name: customer.Name,
	})

    return customer
}
```

### Handling events
To handle events in an aggregate, you need to add an event handler to the aggregate. An event handler is a function that takes a event as an argument and applies it to the aggregate.

```go
// this will add a new event handler for given event name to the aggregate
customer.HandleEvent(name string, handler func(event Event) error)
```

---

## Command
A Command is a request to perform an action in the domain. It is used to change the state of the system.

### Handling a command
Commands are handled by creating a new `command.Bus` and using the `command.Handle` method to subscribe to the command.

```go
// create a new command bus
cmdbus := cqrs.NewInMemoryBus()

// create a new command handler
// note: replace ExampleCommand and ExampleResponse with your own types
commandHandler := func(ctx context.Context, cmd ExampleCommand) (ExampleResponse, error) {
    // handle the command
    // ...
    return nil, nil
}

// handle the command
err := cqrs.Handle[Response, Request any](
    ctx, 
    bus, 
    "command-name", 
    commandHandler,
)
// handle error...
```

### Dispatching a command
To dispatch a command, 
```go
// create a new command bus
bus := cqrs.NewInMemoryBus()

// dispatch the command
// note: replace ExampleCommand and ExampleResponse with your own types
// you can use cqrs.EmptyRequestResponse for commands that don't have response.
resp, err := bus.Dispatch[ExampleResponse](ctx, "command-name", ExampleCommand{})
if err != nil {
	return err
}
```

## Event
The Event represents a event in the domain. It is used to notify other parts of the system that something has happened. An event is represented by the `Event` struct in the `event` package.

### Creating an event
To create an event, you need to create a struct that represents the payload of the event and use the `event.New` function to create a new event.

```go
type CustomerAggregateSampleEventPayload struct {
    CustomField string `json:"custom_field"`
}

func NewCustomAggregateSampleEvent(customerField string) event.Event[CustomerAggregateSampleEventPayload] {
    // create a new event
    ev := event.New(
        "event-id",
        CustomAggregateSampleEventName, 
        CustomerAggregateSampleEventPayload{
            CustomField: customerField,
        },
        // options
        event.WithAggregate(event.Aggregate{
            ID:   "aggregate-id",
            Name: "aggregate-name",
            Version: 1,
        }),
    )

    return ev
}
```

### Handling an event
To subscribe to an event, you need to create a new `event.Bus` and use the `event.Handle` method to subscribe to the event.

```go
// create a new event bus
evbus := event.NewInMemoryBus()

// subscribe to the event
err := event.Handle(ctx, evbus, "event-name", func(ctx event.Context[CustomerAggregateSampleEvent]) error {
    // handle the event
    // ...
    return nil
})
if err != nil {
    return err
}
```

### Publishing an event
To publish an event, you need to create a new `event.Bus` and use the `bus.Publish` method to publish the event.

```go
// create a new event bus
evbus := event.NewInMemoryBus()

// publish the event
err := evbus.Publish(ctx, "event-name", event.Any())
if err != nil {
    return err
}
```


## Repository
A Repository is used to manage the lifecycle of aggregates. It is used to store and retrieve aggregates from a data store.

### Creating an in-memory repository
An in-memory repository is used to store aggregates in memory. It is useful for testing and prototyping.

```go
// create a new in-memory repository
repo := inmemory.NewRepository()

// save the aggregate to the repository
err := repo.Save(ctx, aggregate)
...

// get the aggregate from the repository
agg, err := repo.Get(ctx, "aggregate-id")
...
```