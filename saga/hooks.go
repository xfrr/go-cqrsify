package saga

import "context"

type Hooks struct {
	// OnSagaStarted is called when a saga is started.
	OnSagaStarted func(context.Context, *Instance)
	// OnSagaCompleted is called when a saga completes successfully.
	OnSagaCompleted func(context.Context, *Instance)
	// OnSagaFailed is called when a saga fails and compensation has been triggered.
	OnSagaFailed func(context.Context, *Instance, error)
	// OnSagaCompensating is called when a saga starts compensating.
	OnSagaCompensating func(context.Context, *Instance, int) // from step index
	// OnSagaCompensatingFinished is called when a saga has finished compensating.
	OnSagaCompensatingFinished func(context.Context, *Instance)
	// OnStepStart is called when a step is started.
	OnStepStart func(context.Context, *Instance, StepState)
	// OnStepSuccess is called when a step is completed successfully.
	OnStepSuccess func(context.Context, *Instance, StepState)
	// OnStepFailure is called when a step fails.
	OnStepFailure func(context.Context, *Instance, StepState, error)
	// OnStepCompensationOK is called when a step compensation is successful.
	OnStepCompensationOK func(context.Context, *Instance, StepState)
	// OnStepCompensationKO is called when a step compensation fails.
	OnStepCompensationKO func(context.Context, *Instance, StepState, error)
}
