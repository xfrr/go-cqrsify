package saga_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xfrr/go-cqrsify/messaging"
	messagingmock "github.com/xfrr/go-cqrsify/messaging/mock"
	"github.com/xfrr/go-cqrsify/saga"
)

func newRemoteStepTestExecution() *saga.Execution {
	return &saga.Execution{
		SagaID: "saga-123",
		Def:    &saga.Definition{Steps: []saga.Step{{Name: "reserve-inventory"}}},
		Instance: &saga.Instance{
			Input:    map[string]any{"orderId": "o-1"},
			Metadata: map[string]string{"traceId": "t-1"},
			Steps:    []saga.StepState{{Attempt: 2}},
		},
		StepIndex: 0,
		StepData:  map[string]any{},
	}
}

func TestMessagingRemoteAction_ReturnsContextualErrorWhenRemoteReplyIsNotOK(t *testing.T) {
	t.Parallel()

	ex := newRemoteStepTestExecution()
	bus := &messagingmock.CommandBusReplier{
		PublishRequestFunc: func(_ context.Context, _ messaging.Command) (messaging.Message, error) {
			return saga.RemoteResult{RemoteResultData: saga.RemoteResultData{OK: false, Error: "insufficient funds"}}, nil
		},
		SubscribeWithReplyFunc: func(_ context.Context, _ messaging.MessageHandlerWithReply[messaging.Command, messaging.CommandReply]) (messaging.UnsubscribeFunc, error) {
			return nil, nil
		},
	}

	action := saga.MessagingRemoteAction(bus, saga.RemoteSubjects{Action: "inventory.reserve.cmd"})
	err := action(context.Background(), ex)
	require.Error(t, err)

	var stepErr saga.StepResponseError
	require.True(t, errors.As(err, &stepErr))
	require.Equal(t, "saga-123", stepErr.SagaID)
	require.Equal(t, "reserve-inventory", stepErr.StepName)
	require.Equal(t, "ACTION", stepErr.StepType)
	require.Equal(t, "insufficient funds", stepErr.Reason)

	require.ErrorContains(t, err, "remote ACTION failed")
	require.ErrorContains(t, err, "saga_id=saga-123")
	require.ErrorContains(t, err, "step=reserve-inventory")
	require.ErrorContains(t, err, "insufficient funds")
}

func TestMessagingRemoteCompensation_ReturnsContextualErrorWhenRemoteErrorIsEmpty(t *testing.T) {
	t.Parallel()

	ex := newRemoteStepTestExecution()
	bus := &messagingmock.CommandBusReplier{
		PublishRequestFunc: func(_ context.Context, _ messaging.Command) (messaging.Message, error) {
			return saga.RemoteResult{RemoteResultData: saga.RemoteResultData{OK: false, Error: "   "}}, nil
		},
		SubscribeWithReplyFunc: func(_ context.Context, _ messaging.MessageHandlerWithReply[messaging.Command, messaging.CommandReply]) (messaging.UnsubscribeFunc, error) {
			return nil, nil
		},
	}

	compensate := saga.MessagingRemoteCompensation(bus, saga.RemoteSubjects{Compensate: "inventory.reserve.compensate.cmd"})
	err := compensate(context.Background(), ex)
	require.Error(t, err)

	var stepErr saga.StepResponseError
	require.True(t, errors.As(err, &stepErr))
	require.Equal(t, "saga-123", stepErr.SagaID)
	require.Equal(t, "reserve-inventory", stepErr.StepName)
	require.Equal(t, "COMPENSATE", stepErr.StepType)
	require.Equal(t, "remote replied with ok=false and empty error", stepErr.Reason)

	require.ErrorContains(t, err, "remote COMPENSATE failed")
	require.ErrorContains(t, err, "saga_id=saga-123")
	require.ErrorContains(t, err, "step=reserve-inventory")
	require.ErrorContains(t, err, "remote replied with ok=false and empty error")
}
