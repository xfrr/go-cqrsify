package domain

// EventBusConfig configure an EventBus implementation.
type EventBusConfig struct {
	AsyncWorkers int // >0 enables async worker pool
	QueueSize    int // channel buffer size when async
	// ErrorHandler handles handler failures (after middleware).
	// If nil, errors are logged (if Logger exists) and dropped.
	ErrorHandler func(evtName string, err error)
}

// EventBusConfigModifier is the functional option pattern.
type EventBusConfigModifier func(*EventBusConfig)

func ConfigureEventBusAsyncWorkers(workers int) EventBusConfigModifier {
	return func(o *EventBusConfig) { o.AsyncWorkers = workers }
}

func ConfigureEventBusQueue(size int) EventBusConfigModifier {
	return func(o *EventBusConfig) { o.QueueSize = size }
}

func ConfigureEventBusErrorHandler(fn func(evtName string, err error)) EventBusConfigModifier {
	return func(o *EventBusConfig) { o.ErrorHandler = fn }
}
