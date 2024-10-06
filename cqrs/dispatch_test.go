package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

func TestDispatch(t *testing.T) {
	mockError := errors.New("failed to dispatch command")

	tests := []struct {
		name    string
		ctx     context.Context
		cmd     interface{}
		opts    []cqrs.DispatchOption
		wantErr error
		bus     cqrs.Bus
	}{
		{
			name:    "should fail to dispatch nil bus",
			ctx:     context.Background(),
			cmd:     struct{}{},
			wantErr: cqrs.ErrNilBus,
		},
		{
			name:    "should fail to dispatch nil command",
			ctx:     context.Background(),
			cmd:     nil,
			wantErr: cqrs.ErrInvalidRequest,
			bus:     &mockBus{},
		},
		{
			name:    "should fail to dispatch command",
			ctx:     context.Background(),
			cmd:     struct{}{},
			wantErr: mockError,
			bus: &mockBus{
				dispatch: func(ctx context.Context, cmdname string, cmd interface{}, opts ...cqrs.DispatchOption) (response interface{}, err error) {
					return nil, mockError
				},
			},
		},
		{
			name:    "should dispatch struct command successfully",
			ctx:     context.Background(),
			cmd:     struct{}{},
			wantErr: nil,
			bus: &mockBus{
				dispatch: func(ctx context.Context, cmdname string, cmd interface{}, opts ...cqrs.DispatchOption) (response interface{}, err error) {
					return nil, nil
				},
			},
		},
		{
			name:    "should dispatch string command successfully",
			ctx:     context.Background(),
			cmd:     "Hello, World!",
			wantErr: nil,
			bus:     &mockBus{},
		},
		{
			name:    "should dispatch stringer command successfully",
			ctx:     context.Background(),
			cmd:     cmdStringer{},
			wantErr: nil,
			bus:     &mockBus{},
		},
		{
			name:    "should dispatch gostringer command successfully",
			ctx:     context.Background(),
			cmd:     cmdGoStringer{},
			wantErr: nil,
			bus:     &mockBus{},
		},
		{
			name:    "should dispatch command successfully",
			ctx:     context.Background(),
			cmd:     cmd{},
			wantErr: nil,
			bus:     &mockBus{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := cqrs.Dispatch(tt.ctx, tt.bus, tt.cmd, tt.opts...)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error to be %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil && res != nil {
				t.Fatalf("expected response to be nil, got %v", res)
			}
		})
	}
}
