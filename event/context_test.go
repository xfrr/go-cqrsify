package event_test

import (
	"context"
	"testing"

	"github.com/xfrr/cqrsify/event"
)

func TestContext(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		type args struct {
			ctx context.Context
			evt event.Event[any]
		}

		cases := []struct {
			name string
			args args
		}{
			{
				name: "should create a new context",
				args: args{
					ctx: context.Background(),
					evt: event.New("id", "name", "msg").Any(),
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				ctx := event.WithContext(tt.args.ctx, tt.args.evt)
				if ctx.Event().ID() != tt.args.evt.ID() {
					t.Errorf("expected event id to be %v, got %v", tt.args.evt.ID(), ctx.Event().ID())
				}
				if ctx.Event().Payload() != tt.args.evt.Payload() {
					t.Errorf("expected event message to be %v, got %v", tt.args.evt.Payload(), ctx.Event().Payload())
				}
				if ctx.Event().Reason() != tt.args.evt.Reason() {
					t.Errorf("expected event reason to be %v, got %v", tt.args.evt.Reason(), ctx.Event().Reason())
				}
				if ctx.Event().Time() != tt.args.evt.Time() {
					t.Errorf("expected event time to be %v, got %v", tt.args.evt.Time(), ctx.Event().Time())
				}
			})
		}
	})

	t.Run("CastContext", func(t *testing.T) {
		t.Run("should return false when cast fails", func(t *testing.T) {
			ctx := event.WithContext(context.Background(), event.New("id", "name", 1))
			_, ok := event.CastContext[string](ctx)
			if ok {
				t.Errorf("expected cast to fail")
			}
		})

		t.Run("should return true when cast succeeds", func(t *testing.T) {
			casted, ok := event.CastContext[string](event.WithContext(context.Background(), event.New("id", "name", "msg")))
			if !ok {
				t.Errorf("expected cast to succeed")
			}
			if casted.Event().ID() != "id" {
				t.Errorf("expected casted event id to be %v, got %v", "id", casted.Event().ID())
			}
			if casted.Event().Payload() != "msg" {
				t.Errorf("expected casted event message to be %v, got %v", "msg", casted.Event().Payload())
			}
		})
	})
}
