package message_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/xfrr/go-cqrsify/message"
)

// Test message types
type TestMessage struct {
	message.Base
}

type AnotherTestMessage struct {
	message.Base
}

// Mock handler for testing
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(ctx context.Context, msg message.Message) (any, error) {
	args := m.Called(ctx, msg)
	switch res := args.Get(0).(type) {
	case string:
		return res, nil
	case error:
		return nil, res
	}

	return nil, nil
}

// Test suite structure
type InMemoryBusTestSuite struct {
	suite.Suite
	bus *message.InMemoryBus
}

func (suite *InMemoryBusTestSuite) SetupTest() {
	suite.bus = message.NewInMemoryBus()
}

// Test NewInMemoryBus constructor
func (suite *InMemoryBusTestSuite) TestNewInMemoryBus() {
	bus := message.NewInMemoryBus()

	assert.NotNil(suite.T(), bus)
}

// Test successful handler registration
func (suite *InMemoryBusTestSuite) TestRegisterHandler_Success() {
	handler := &MockHandler{}

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	assert.NoError(suite.T(), err)
}

// Test duplicate handler registration error
func (suite *InMemoryBusTestSuite) TestRegisterHandler_DuplicateError() {
	handler1 := &MockHandler{}
	handler2 := &MockHandler{}

	err1 := suite.bus.RegisterHandler("com.org.test_message", handler1)
	require.NoError(suite.T(), err1)

	err2 := suite.bus.RegisterHandler("com.org.test_message", handler2)

	assert.Error(suite.T(), err2)
	assert.Contains(suite.T(), err2.Error(), "handler already registered for message type com.org.test_message")
}

// Test successful message dispatch
func (suite *InMemoryBusTestSuite) TestDispatch_Success() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	handler.
		On("Handle", ctx, msg).
		Return("ok")

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", res)
	handler.AssertExpectations(suite.T())
}

// Test dispatch with handler error
func (suite *InMemoryBusTestSuite) TestDispatch_HandlerError() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()
	expectedError := errors.New("handler error")

	handler.On("Handle", ctx, msg).Return(expectedError)

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), expectedError, err)
	handler.AssertExpectations(suite.T())
}

// Test dispatch with no registered handler
func (suite *InMemoryBusTestSuite) TestDispatch_NoHandlerError() {
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Contains(suite.T(), err.Error(), "no handler registered for message type com.org.test_message")
}

// Test dispatch with context cancellation
func (suite *InMemoryBusTestSuite) TestDispatch_WithCancelledContext() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	handler.On("Handle", ctx, msg).Return(context.Canceled)

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), context.Canceled, err)
	handler.AssertExpectations(suite.T())
}

// Test dispatch with timeout context
func (suite *InMemoryBusTestSuite) TestDispatch_WithTimeout() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	handler.On("Handle", ctx, msg).Return(nil)

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
	handler.AssertExpectations(suite.T())
}

// Test middleware application
func (suite *InMemoryBusTestSuite) TestUse_MiddlewareApplication() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	middlewareCalled := false
	middleware := func(h message.Handler[message.Message, any]) message.Handler[message.Message, any] {
		return &handlerWrapper{
			fn: func(ctx context.Context, msg message.Message) (any, error) {
				middlewareCalled = true
				return h.Handle(ctx, msg)
			},
		}
	}

	handler.On("Handle", ctx, msg).Return(nil)

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	suite.bus.Use(middleware)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.True(suite.T(), middlewareCalled)
	handler.AssertExpectations(suite.T())
}

// Test multiple middlewares in correct order
func (suite *InMemoryBusTestSuite) TestUse_MultipleMiddlewares() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	var executionOrder []string

	middleware1 := func(h message.Handler[message.Message, any]) message.Handler[message.Message, any] {
		return &handlerWrapper{
			fn: func(ctx context.Context, msg message.Message) (any, error) {
				executionOrder = append(executionOrder, "middleware1_before")
				res, err := h.Handle(ctx, msg)
				executionOrder = append(executionOrder, "middleware1_after")
				return res, err
			},
		}
	}

	middleware2 := func(h message.Handler[message.Message, any]) message.Handler[message.Message, any] {
		return &handlerWrapper{
			fn: func(ctx context.Context, msg message.Message) (any, error) {
				executionOrder = append(executionOrder, "middleware2_before")
				res, err := h.Handle(ctx, msg)
				executionOrder = append(executionOrder, "middleware2_after")
				return res, err
			},
		}
	}

	handler.On("Handle", ctx, msg).Run(func(args mock.Arguments) {
		executionOrder = append(executionOrder, "handler")
	}).Return(nil)

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	// Add middlewares first, then use Use() to apply them
	suite.bus.Use(middleware1)
	suite.bus.Use(middleware2)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)

	// Middlewares should be applied in reverse order (last added, first executed)
	expectedOrder := []string{
		"middleware2_before",
		"middleware1_before",
		"handler",
		"middleware1_after",
		"middleware2_after",
	}
	assert.Equal(suite.T(), expectedOrder, executionOrder)
	handler.AssertExpectations(suite.T())
}

// Test middleware with error
func (suite *InMemoryBusTestSuite) TestUse_MiddlewareWithError() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()
	expectedError := errors.New("middleware error")

	middleware := func(h message.Handler[message.Message, any]) message.Handler[message.Message, any] {
		return &handlerWrapper{
			fn: func(ctx context.Context, msg message.Message) (any, error) {
				return nil, expectedError
			},
		}
	}

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	suite.bus.Use(middleware)

	res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), expectedError, err)
	// Handler should not be called when middleware returns error
	handler.AssertNotCalled(suite.T(), "Handle")
}

// Test concurrent access safety
func (suite *InMemoryBusTestSuite) TestConcurrentAccess() {
	handler := &MockHandler{}
	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	handler.
		On("Handle", mock.Anything, mock.Anything).
		Return("ok")

	err := suite.bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(suite.T(), err)

	// Run multiple goroutines concurrently
	const numGoroutines = 10
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			res, err := suite.bus.Dispatch(ctx, "com.org.test_message", msg)
			assert.Equal(suite.T(), "ok", res)
			errChan <- err
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-errChan
		assert.NoError(suite.T(), err)
	}
}

// Test different message types
func (suite *InMemoryBusTestSuite) TestMultipleMessageTypes() {
	handler1 := &MockHandler{}
	handler2 := &MockHandler{}

	msg1 := TestMessage{message.NewBase()}
	msg2 := AnotherTestMessage{message.NewBase()}
	ctx := context.Background()

	handler1.On("Handle", ctx, msg1).Return(nil)
	handler2.On("Handle", ctx, msg2).Return(nil)

	err1 := suite.bus.RegisterHandler("com.org.test_message", handler1)
	err2 := suite.bus.RegisterHandler("com.org.another_test_message", handler2)
	require.NoError(suite.T(), err1)
	require.NoError(suite.T(), err2)

	res1, err1 := suite.bus.Dispatch(ctx, "com.org.test_message", msg1)
	res2, err2 := suite.bus.Dispatch(ctx, "com.org.another_test_message", msg2)
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.Nil(suite.T(), res1)
	assert.Nil(suite.T(), res2)
	handler1.AssertExpectations(suite.T())
	handler2.AssertExpectations(suite.T())
}

// Test handlerWrapper implementation (helper for tests)
type handlerWrapper struct {
	fn func(ctx context.Context, msg message.Message) (any, error)
}

func (h *handlerWrapper) Handle(ctx context.Context, msg message.Message) (any, error) {
	return h.fn(ctx, msg)
}

// Test Bus interface compliance
func (suite *InMemoryBusTestSuite) TestBusInterfaceCompliance() {
	var bus message.Bus = suite.bus
	assert.NotNil(suite.T(), bus)
}

// Run the test suite
func TestInMemoryBusTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryBusTestSuite))
}

// Additional table-driven tests for edge cases
func TestInMemoryBus_EdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		setupFunc       func(*message.InMemoryBus) error
		message         message.Message
		expectedError   string
		shouldHaveError bool
	}{
		{
			name: "dispatch nil message",
			setupFunc: func(bus *message.InMemoryBus) error {
				handler := &MockHandler{}
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil)
				return bus.RegisterHandler("com.org.test_message", handler)
			},
			message:         nil,
			expectedError:   "no handler registered for message type",
			shouldHaveError: true,
		},
		{
			name: "register handler with empty message type",
			setupFunc: func(bus *message.InMemoryBus) error {
				handler := &MockHandler{}
				return bus.RegisterHandler("", handler)
			},
			message:         TestMessage{},
			expectedError:   "no handler registered for message type com.org.test_message",
			shouldHaveError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := message.NewInMemoryBus()

			if tt.setupFunc != nil {
				setupErr := tt.setupFunc(bus)
				if tt.name == "register handler with empty message type" {
					// We expect this setup to succeed, but dispatch to fail
					assert.NoError(t, setupErr)
				}
			}

			if tt.message != nil {
				res, err := bus.Dispatch(context.Background(), "com.org.test_message", tt.message)
				if tt.shouldHaveError {
					assert.Nil(t, res)
					assert.Error(t, err)
					if tt.expectedError != "" {
						assert.Contains(t, err.Error(), tt.expectedError)
					}
				} else {
					assert.Nil(t, res)
					assert.NoError(t, err)
				}
			}
		})
	}
}
