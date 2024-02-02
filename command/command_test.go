package command_test

import (
	"testing"

	"github.com/xfrr/cqrsify/command"
)

func TestCommand(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		type args struct {
			id            string
			name          string
			msg           string
			aggregateName string
			aggregateID   string
		}

		cases := []struct {
			name string
			args args
		}{
			{
				name: "with aggregate",
				args: args{
					id:            "id",
					name:          "name",
					msg:           "msg",
					aggregateName: "aggregateName",
					aggregateID:   "aggregateID",
				},
			},
			{
				name: "without aggregate",
				args: args{
					id:            "id",
					name:          "name",
					msg:           "msg",
					aggregateName: "",
					aggregateID:   "",
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				cmd := command.New(tt.args.id, tt.args.msg, command.WithAggregate(tt.args.aggregateID, tt.args.aggregateName))
				if cmd.ID() != command.ID(tt.args.id) {
					t.Errorf("ID() = %v, want %v", cmd.ID(), tt.args.id)
				}
				if cmd.Message() != tt.args.msg {
					t.Errorf("Message() = %v, want %v", cmd.Message(), tt.args.msg)
				}
				if cmd.AggregateName() != tt.args.aggregateName {
					t.Errorf("AggregateName() = %v, want %v", cmd.AggregateName(), tt.args.aggregateName)
				}
				if cmd.AggregateID() != tt.args.aggregateID {
					t.Errorf("AggregateID() = %v, want %v", cmd.AggregateID(), tt.args.aggregateID)
				}
			})
		}
	})

	t.Run("Any", func(t *testing.T) {
		cmd := command.New("id", "msg")
		any := cmd.Any()
		if any.ID() != cmd.ID() {
			t.Errorf("ID() = %v, want %v", any.ID(), cmd.ID())
		}
		if any.Message() != cmd.Message() {
			t.Errorf("Message() = %v, want %v", any.Message(), cmd.Message())
		}
		if any.AggregateName() != cmd.AggregateName() {
			t.Errorf("AggregateName() = %v, want %v", any.AggregateName(), cmd.AggregateName())
		}
		if any.AggregateID() != cmd.AggregateID() {
			t.Errorf("AggregateID() = %v, want %v", any.AggregateID(), cmd.AggregateID())
		}
	})

	t.Run("Cast", func(t *testing.T) {
		t.Run("should return false when cast fails", func(t *testing.T) {
			_, ok := command.Cast[int](command.New("id", "msg"))
			if ok {
				t.Errorf("Cast() = %v, want %v", ok, false)
			}
		})

		t.Run("should return true when cast succeeds", func(t *testing.T) {
			casted, ok := command.Cast[string](command.New("id", "msg"))
			if !ok {
				t.Errorf("Cast() = %v, want %v", ok, true)
			}

			if casted.ID() != "id" {
				t.Errorf("ID() = %v, want %v", casted.ID(), "id")
			}
			if casted.Message() != "msg" {
				t.Errorf("Message() = %v, want %v", casted.Message(), "msg")
			}
		})

		t.Run("should return true when cast succeeds with aggregate", func(t *testing.T) {
			casted, ok := command.Cast[string](command.New("id", "msg", command.WithAggregate("aggregateID", "aggregateName")))
			if !ok {
				t.Errorf("Cast() = %v, want %v", ok, true)
			}

			if casted.AggregateName() != "aggregateName" {
				t.Errorf("AggregateName() = %v, want %v", casted.AggregateName(), "aggregateName")
			}
			if casted.AggregateID() != "aggregateID" {
				t.Errorf("AggregateID() = %v, want %v", casted.AggregateID(), "aggregateID")
			}
		})
	})
}
