package command_test

import (
	"context"
	"testing"

	"github.com/xfrr/cqrsify/command"
)

func TestContext(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		type args struct {
			ctx context.Context
			cmd command.Command[any]
		}

		cases := []struct {
			name string
			args args
		}{
			{
				name: "with aggregate",
				args: args{
					ctx: context.Background(),
					cmd: command.New("id", "payload",
						command.WithAggregate("aggregateID", "aggregateName"),
					).Any(),
				},
			},
			{
				name: "without aggregate",
				args: args{
					ctx: context.Background(),
					cmd: command.New("id", "payload").Any(),
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				ctx := command.WithContext(tt.args.ctx, tt.args.cmd)
				if ctx.Command().ID() != tt.args.cmd.ID() {
					t.Errorf("expected command id to be %v, got %v", tt.args.cmd.ID(), ctx.Command().ID())
				}
				if ctx.Command().Payload() != tt.args.cmd.Payload() {
					t.Errorf("expected command payload to be %v, got %v", tt.args.cmd.Payload(), ctx.Command().Payload())
				}
				if ctx.Command().AggregateName() != tt.args.cmd.AggregateName() {
					t.Errorf("expected command aggregate name to be %v, got %v", tt.args.cmd.AggregateName(), ctx.Command().AggregateName())
				}
				if ctx.Command().AggregateID() != tt.args.cmd.AggregateID() {
					t.Errorf("expected command aggregate id to be %v, got %v", tt.args.cmd.AggregateID(), ctx.Command().AggregateID())
				}
			})
		}
	})

	t.Run("CastContext", func(t *testing.T) {
		t.Run("should return false when cast fails", func(t *testing.T) {
			ctx := command.WithContext(context.Background(), command.New("id", 1))
			_, ok := command.CastContext[string](ctx)
			if ok {
				t.Errorf("expected cast to fail")
			}
		})

		t.Run("should return true when cast succeeds", func(t *testing.T) {
			casted, ok := command.CastContext[string](command.WithContext(context.Background(), command.New("id", "payload")))
			if !ok {
				t.Errorf("expected cast to succeed")
			}
			if casted.Command().ID() != "id" {
				t.Errorf("expected casted command id to be %v, got %v", "id", casted.Command().ID())
			}
			if casted.Command().Payload() != "payload" {
				t.Errorf("expected casted command payload to be %v, got %v", "payload", casted.Command().Payload())
			}
		})

		t.Run("should return true when cast succeeds with aggregate", func(t *testing.T) {
			casted, ok := command.CastContext[string](command.WithContext(context.Background(), command.New("id", "payload",
				command.WithAggregate("aggregateID", "aggregateName"),
			)))
			if !ok {
				t.Errorf("expected cast to succeed")
			}

			if casted.Command().AggregateName() != "aggregateName" {
				t.Errorf("expected casted command aggregate name to be %v, got %v", "aggregateName", casted.Command().AggregateName())
			}
			if casted.Command().AggregateID() != "aggregateID" {
				t.Errorf("expected casted command aggregate id to be %v, got %v", "aggregateID", casted.Command().AggregateID())
			}
		})
	})
}
