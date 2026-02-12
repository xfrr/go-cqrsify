package messaging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xfrr/go-cqrsify/messaging"
)

type InMemoryMessageBusTestSuite struct {
	suite.Suite

	sut *messaging.InMemoryMessageBus
}

func TestInMemoryMessageBusSuite(t *testing.T) {
	suite.Run(t, new(InMemoryMessageBusTestSuite))
}

func (s *InMemoryMessageBusTestSuite) SetupTest() {
	s.sut = messaging.NewInMemoryMessageBus()
}

func (s *InMemoryMessageBusTestSuite) TearDownTest() {
	err := s.sut.Close()
	s.Require().NoError(err)
}

func (s *InMemoryMessageBusTestSuite) TestNewMessageBus() {
	s.Require().NotNil(s.sut)
}

func (s *InMemoryMessageBusTestSuite) TestPublish() {
	const subject = "test.message.type.publish"
	bus := messaging.NewInMemoryMessageBus(
		messaging.ConfigureInMemoryMessageBusSubjects(subject),
	)

	msg := messaging.NewMessage(subject)

	s.Run("should succeed when handler is subscribed", func() {
		seen := make(chan messaging.Message, 1)

		unsub, err := bus.Subscribe(
			s.T().Context(),
			messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, e messaging.Message) error {
				seen <- e
				return nil
			}),
		)
		s.Require().NoError(err)

		err = bus.Publish(s.T().Context(), msg)
		s.Require().NoError(err)

		select {
		case got := <-seen:
			s.Require().Equal(subject, got.MessageType())
		case <-s.T().Context().Done():
			s.T().Fatalf("handler was not invoked for %q", subject)
		}

		err = unsub()
		s.Require().NoError(err)
	})

	s.Run("should return error when no handler is subscribed", func() {
		otherMsg := messaging.NewMessage("other.message.type")
		err := bus.Publish(s.T().Context(), otherMsg)
		s.Require().ErrorIs(err, messaging.NoHandlersForMessageError{MessageType: otherMsg.MessageType()})
	})

	s.Run("should return error when empty message list is published", func() {
		err := bus.Publish(s.T().Context())
		s.Require().ErrorContains(err, "no messages to publish")
	})

	s.Run("should return error when handler returns error", func() {
		expectedErr := messaging.InvalidMessageTypeError{Expected: "expected", Actual: "actual"}
		unsub, err := bus.Subscribe(
			s.T().Context(),
			messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
				return expectedErr
			}),
		)
		s.Require().NoError(err)

		err = bus.Publish(s.T().Context(), msg)
		s.Require().ErrorIs(err, expectedErr)

		err = unsub()
		s.Require().NoError(err)
	})

	s.Run("should return error when publisher is closed", func() {
		err := bus.Close()
		s.Require().NoError(err)

		err = bus.Publish(s.T().Context(), msg)
		s.Require().ErrorIs(err, messaging.ErrPublishOnClosedBus)
	})
}

func (s *InMemoryMessageBusTestSuite) TestPublishRequest() {
	const subject = "test.message.type.request"
	bus := messaging.NewInMemoryMessageBus(
		messaging.ConfigureInMemoryMessageBusSubjects(subject),
	)

	msg := messaging.NewMessage(subject)

	s.Run("should succeed when handler is subscribed", func() {
		expectedReply := messaging.NewMessage("reply.message.type")
		unsub, err := bus.SubscribeWithReply(
			s.T().Context(),
			messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.Message](func(_ context.Context, _ messaging.Message) (messaging.Message, error) {
				return expectedReply, nil
			}),
		)
		s.Require().NoError(err)

		reply, err := bus.PublishRequest(s.T().Context(), msg)
		s.Require().NoError(err)
		s.Require().Equal(expectedReply, reply)

		err = unsub()
		s.Require().NoError(err)
	})

	s.Run("should return error when no handler is subscribed", func() {
		otherMsg := messaging.NewMessage("other.message.type")
		unsub, err := bus.PublishRequest(s.T().Context(), otherMsg)
		s.Require().ErrorIs(err, messaging.NoHandlersForMessageError{MessageType: otherMsg.MessageType()})
		s.Require().Nil(unsub)
	})

	s.Run("should return error when handler returns error", func() {
		expectedErr := messaging.InvalidMessageTypeError{Expected: "expected", Actual: "actual"}
		unsub, err := bus.SubscribeWithReply(
			s.T().Context(),
			messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.Message](func(_ context.Context, _ messaging.Message) (messaging.Message, error) {
				return nil, expectedErr
			}),
		)
		s.Require().NoError(err)

		_, err = bus.PublishRequest(s.T().Context(), msg)
		s.Require().ErrorIs(err, expectedErr)

		err = unsub()
		s.Require().NoError(err)
	})
}

func (s *InMemoryMessageBusTestSuite) TestNewMessageBus_SubscribeAndPublish_Worker_Suceeds() {
	const subject = "test.message.type.3"
	bus := messaging.NewInMemoryMessageBus(
		messaging.ConfigureInMemoryMessageBusSubjects(subject),
		messaging.ConfigureInMemoryMessageBusAsyncWorkers(2),
		messaging.ConfigureInMemoryMessageBusQueueBufferSize(2),
	)

	msg1 := messaging.NewMessage(subject)
	msg2 := messaging.NewMessage(subject)

	seen := make(chan messaging.Message, 1)

	unsub, err := bus.Subscribe(
		s.T().Context(),
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, e messaging.Message) error {
			seen <- e
			return nil
		}),
	)
	s.Require().NoError(err)

	err = bus.Publish(s.T().Context(), msg1, msg2)
	s.Require().NoError(err)

	for range 2 {
		select {
		case got := <-seen:
			s.Require().Equal(subject, got.MessageType())
		case <-s.T().Context().Done():
			s.T().Fatalf("handler was not invoked for %q", subject)
		}
	}

	err = unsub()
	s.Require().NoError(err)
}

func (s *InMemoryMessageBusTestSuite) TestNewMessageBus_SubscribeAndPublish_Worker_Handler_ReturnsError() {
	const subject = "test.message.type.5"

	errCh := make(chan error, 1)
	errHandler := func(messageType string, err error) {
		s.Require().Equal(subject, messageType)
		errCh <- err
	}

	bus := messaging.NewInMemoryMessageBus(
		messaging.ConfigureInMemoryMessageBusSubjects(subject),
		messaging.ConfigureInMemoryMessageBusAsyncWorkers(1),
		messaging.ConfigureInMemoryMessageBusQueueBufferSize(1),
		messaging.ConfigureInMemoryMessageBusErrorHandler(errHandler),
	)

	msg1 := messaging.NewMessage(subject)

	expectedErr := messaging.InvalidMessageTypeError{Expected: "expected", Actual: "actual"}
	unsub, err := bus.Subscribe(
		s.T().Context(),
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
			return expectedErr
		}),
	)
	s.Require().NoError(err)

	err = bus.Publish(s.T().Context(), msg1)
	s.Require().NoError(err)

	select {
	case gotErr := <-errCh:
		s.Require().ErrorIs(gotErr, expectedErr)
	case <-s.T().Context().Done():
		s.T().Fatalf("error handler was not invoked for %q", subject)
	}

	err = unsub()
	s.Require().NoError(err)
}
