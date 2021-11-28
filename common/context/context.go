package context

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// Context application context wrapping base context
type Context struct {
	context.Context
}

func FromContext(ctx context.Context) Context {
	switch c := ctx.(type) {
	case Context:
		return c
	default:
		return Context{Context: ctx}
	}
}

// GetActorContext returns additional actor.Context from application context
func (Context) GetActorContext() actor.Context {
	return nil
}
