package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestContext(t *testing.T) {
	t.Run("WithContext", func(t *testing.T) {
		type args struct {
			ctx context.Context
			evt func() event.Event[any, any]
		}

		cases := []struct {
			name string
			args args
		}{
			{
				name: "should create a new context",
				args: args{
					ctx: context.Background(),
					evt: func() event.Event[any, any] {
						evt, _ := event.New(
							"id",
							"name",
							"payload",
							event.WithOccurredAt(
								time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							))
						return evt.Any()
					},
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				ctx := event.WithContext(tt.args.ctx, tt.args.evt())
				if ctx.Event().ID() != tt.args.evt().ID() {
					t.Errorf("expected event id to be %v, got %v", tt.args.evt().ID(), ctx.Event().ID())
				}
				if ctx.Event().Payload() != tt.args.evt().Payload() {
					t.Errorf("expected event payload to be %v, got %v", tt.args.evt().Payload(), ctx.Event().Payload())
				}
				if ctx.Event().Name() != tt.args.evt().Name() {
					t.Errorf("expected event name to be %v, got %v", tt.args.evt().Name(), ctx.Event().Name())
				}
				if ctx.Event().OccurredAt() != tt.args.evt().OccurredAt() {
					t.Errorf("expected event time to be %v, got %v", tt.args.evt().OccurredAt(), ctx.Event().OccurredAt())
				}
			})
		}
	})

	t.Run("CastContext", func(t *testing.T) {
		t.Run("should return false when cast fails", func(t *testing.T) {
			evt, err := event.New("id", "name", 1)
			if err != nil {
				t.Fatal(err)
			}

			ctx := event.WithContext(context.Background(), evt)
			_, ok := event.CastContext[int, any](ctx)
			if ok {
				t.Errorf("expected cast to fail")
			}
		})

		t.Run("should return true when cast succeeds", func(t *testing.T) {
			evt, err := event.New("id", "name", "payload")
			if err != nil {
				t.Fatal(err)
			}
			casted, ok := event.CastContext[string, string, string, string](event.WithContext(context.Background(), evt))
			if !ok {
				t.Errorf("expected cast to succeed")
			}
			if casted.Event().ID() != "id" {
				t.Errorf("expected casted event id to be %v, got %v", "id", casted.Event().ID())
			}
			if casted.Event().Payload() != "payload" {
				t.Errorf("expected casted event payload to be %v, got %v", "payload", casted.Event().Payload())
			}
		})
	})
}
