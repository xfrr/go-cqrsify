package messaging

import "context"

// PublishEvent is a shorthand for publishing events.
func PublishEvent(ctx context.Context, publisher EventPublisher, event Event) error {
	err := publisher.Publish(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
