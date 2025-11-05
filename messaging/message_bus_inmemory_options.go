package messaging

// MessageBusConfig configure an MessageBus implementation.
type MessageBusConfig struct {
	AsyncWorkers int // >0 enables async worker pool
	QueueSize    int // channel buffer size when async
	// ErrorHandler handles handler failures (after middleware).
	// If nil, errors are logged (if Logger exists) and dropped.
	ErrorHandler func(evtName string, err error)
	// Subjects is a list of subjects the bus listens to. If empty, subscribes to all messages.
	Subjects []string
}

// MessageBusConfigConfiger is the functional option pattern.
type MessageBusConfigConfiger func(*MessageBusConfig)

func ConfigureInMemoryMessageBusAsyncWorkers(workers int) MessageBusConfigConfiger {
	return func(o *MessageBusConfig) { o.AsyncWorkers = workers }
}

func ConfigureInMemoryMessageBusQueueBufferSize(size int) MessageBusConfigConfiger {
	return func(o *MessageBusConfig) { o.QueueSize = size }
}

func ConfigureInMemoryMessageBusErrorHandler(fn func(evtName string, err error)) MessageBusConfigConfiger {
	return func(o *MessageBusConfig) { o.ErrorHandler = fn }
}

func ConfigureInMemoryMessageBusSubjects(subjects ...string) MessageBusConfigConfiger {
	return func(o *MessageBusConfig) { o.Subjects = subjects }
}
