package cqrs_test

import (
	"context"
	"sync"

	"github.com/xfrr/go-cqrsify/cqrs"
)

type dispatchCall struct {
	ctx     context.Context
	cmdname string
	cmd     interface{}
	opts    []cqrs.DispatchOption
}

type registerCall struct {
	ctx     context.Context
	cmdname string
	handler cqrs.HandlerFuncAny
}

type mockBus struct {
	lock          sync.Mutex
	dispatchCalls []dispatchCall
	dispatch      func(ctx context.Context, cmdname string, cmd interface{}, opts ...cqrs.DispatchOption) (response interface{}, err error)
	registerCalls []registerCall
	register      func(ctx context.Context, cmdname string, handler cqrs.HandlerFuncAny) error
}

func (m *mockBus) Dispatch(ctx context.Context, cmdname string, cmd interface{}, opts ...cqrs.DispatchOption) (response interface{}, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.dispatchCalls = append(m.dispatchCalls, dispatchCall{ctx, cmdname, cmd, opts})
	if m.dispatch != nil {
		return m.dispatch(ctx, cmdname, cmd, opts...)
	}
	return nil, nil
}

func (m *mockBus) RegisterHandler(ctx context.Context, cmdname string, handler cqrs.HandlerFuncAny) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.registerCalls = append(m.registerCalls, registerCall{ctx, cmdname, handler})
	if m.register != nil {
		return m.register(ctx, cmdname, handler)
	}
	return nil
}
