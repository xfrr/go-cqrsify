package messaging

// MessageBusConfig configure an MessageBus implementation.
type MessageBusConfig struct {
	AsyncWorkers int // >0 enables async worker pool
	QueueSize    int // channel buffer size when async
	// ErrorHandler handles handler failures (after middleware).
	// If nil, errors are logged (if Logger exists) and dropped.
	ErrorHandler func(evtName string, err error)
}

// MessageBusConfigModifier is the functional option pattern.
type MessageBusConfigModifier func(*MessageBusConfig)

func ConfigureMessageBusAsyncWorkers(workers int) MessageBusConfigModifier {
	return func(o *MessageBusConfig) { o.AsyncWorkers = workers }
}

func ConfigureMessageBusQueue(size int) MessageBusConfigModifier {
	return func(o *MessageBusConfig) { o.QueueSize = size }
}

func ConfigureMessageBusErrorHandler(fn func(evtName string, err error)) MessageBusConfigModifier {
	return func(o *MessageBusConfig) { o.ErrorHandler = fn }
}
