package aggregate

import (
	"github.com/xfrr/go-cqrsify/aggregate/event"
)

// Event represents an event that events the state of an Aggregate.
// It is an alias for event.Event[any].
type Event = event.Event[any, any]
