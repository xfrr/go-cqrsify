
# Usage
This section explains how to use Cqrsify in your application. You can find examples of how to use the library in the [examples](https://github.com/xfrr/cqrsify/tree/main/examples).

## Aggregate
An aggregate is a domain-driven design (DDD) concept that represents a cluster of domain objects that can be treated as a single unit. An aggregate will have a unique identifier and a set of events that represent changes to the aggregate.

> **Note**: An Event is called a Change in the Aggregate context and is represented by the `Change` struct in the `aggregate` package.

### Creating an aggregate
To create an aggregate, you need to create a struct that embeds the `*aggregate.Base` struct and implements the `aggregate.Aggregate` interface. The `Base` struct provides the basic functionality for an aggregate, such as the `ApplyChange` method, which is used to apply a change to the aggregate.

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

    // apply the change to the aggregate
    // note: this is a helper function that creates an event and applies it to the 
    // aggregate. You can also create the event manually and apply it to the aggregate 
    // by calling customer.ApplyChange(event)
	aggregate.ApplyChange(
		customer,
		CustomAggregateCreatedEventName,
		CustomAggregateCreatedEvent{
			ID:   customer.ID().String(),
			Name: customer.Name,
	})

    // some other aggregate methods
    // customer.ID() // returns the aggregate id
    // customer.Name() // returns the aggregate name

    return customer
}
```

### Handling changes
To handle changes in an aggregate, you need to add an event handler to the aggregate. An event handler is a function that takes a change as an argument. The `When` method is used to add a new event handler to the aggregate.

```go
// this will add a new event handler for given event reason to the aggregate
customer.When(reason string, handler func(change Change) error)
```

### Applying changes
To apply a change to the aggregate, you need to create an event and apply it to the aggregate. The `ApplyChange` method is used to apply a change to the aggregate.

```go
// this will apply the change (event) to the aggregate
customer.ApplyChange(change Change)
```

---

## Command
A Command is a request to perform an action. It is used to represent the intention of the user to change the state of the system. A command is represented by the `Command` struct in the `command` package.

### Creating a command
To create a command, you need to create a struct that represents the payload of the command and use the `command.New` function to create a new command.

```go
type CustomerAggregateSampleCommand struct {
    CustomField string `json:"custom_field"`
}

func NewCustomerAggregateSampleCommand(commandID, customField string) command.Command[CustomerAggregateSampleCommandPayload] {
    // create a new command
   	cmd := command.New[CustomerAggregateSampleCommand](
        commandID,
        CustomerAggregateSampleCommand{
            CustomField: customField,
        },
		command.WithAggregate(MockAggregateID, MockAggregateName),
	)

    return cmd
}
```

### Handling a command
To handle a command, you need to create a new `command.Bus` and creates a new `command.Handler` to handle the command.

```go
// create a new command bus
cmdbus := command.NewBus()

// this will subscribe to the command and handle it asynchronously.
// it returns a chan of errors that occurred when handling the command
errs, err := command.Handle(ctx, bus, "command-name", func(ctx command.Context[CustomerAggregateSampleCommand]) error {
    // handle the command
    // ...
    return nil
})
if err != nil {
	return err
}

// listen for errors
go func() {
    for err := range errs {
        // handle the error
    }
}()
```

### Dispatching a command
To dispatch a command, you need to create a new `command.Bus` and use the `command.Dispatch` method to dispatch the command.

```go
// create a new command bus
cmdbus := command.NewBus()

// dispatch the command
err := bus.Dispatch(ctx, "command-name", command.Any())
if err != nil {
	return err
}
```

## Event
The Event represents a change in the domain. It is used to notify other parts of the system that something has happened. An event is represented by the `Event` struct in the `event` package.

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
        CustomAggregateSampleEventReason, 
        CustomerAggregateSampleEventPayload{
            CustomField: customerField,
        },
        // options
        event.WithAggregate(event.Aggregate{
            ID:   "aggregate-id",
            Name: "aggregate-name",
            Version: 1,
        }),
        event.WithVersion(1),
    )

    return ev
}
```

### Handling an event
To subscribe to an event, you need to create a new `event.Bus` and use the `event.Handle` method to subscribe to the event.

```go
// create a new event bus
evbus := event.NewBus()

// subscribe to the event
err := event.Handle(ctx, evbus, "event-reason", func(ctx event.Context[CustomerAggregateSampleEvent]) error {
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
evbus := event.NewBus()

// publish the event
err := evbus.Publish(ctx, "event-reason", event.Any())
if err != nil {
    return err
}
```

## Policy
A Policy is used to enforce business rules. It is used to validate the state of the system and ensure that the system is in a consistent state.

### Creating a policy
*Coming Soon...*

### Enforcing a policy
*Coming Soon...*



## Query
A Query is a request to retrieve data from the system.

### Creating a query
*Coming Soon...*

## Repository
A Repository is used to manage the lifecycle of aggregates. It is used to store and retrieve aggregates from a data store.

### Creating a repository
*Coming Soon...*

## Snapshot
A Snapshot is a point-in-time representation of an aggregate. It is used to optimize the performance of the system by reducing the number of events that need to be replayed.

### Creating a snapshot
*Coming Soon...*

### Restoring a snapshot
*Coming Soon...*
