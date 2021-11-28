package context

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// Context application context wrapping base context
type Context interface {
	context.Context
	GetActorContext() actor.Context
}

type key struct{ name string }

var actorCtxKey = key{"actor_ctx_key"} // private context key is important

type appContext struct {
	context.Context
}

func Background() Context {
	return &appContext{context.Background()}
}

func TODO() Context {
	return &appContext{context.TODO()}
}

func FromContext(ctx context.Context) Context {
	switch c := ctx.(type) {
	case *appContext:
		return c
	default:
		return &appContext{Context: ctx}
	}
}

func (c *appContext) WithActorContext(actx actor.Context) Context {
	c.Context = context.WithValue(c.Context, actorCtxKey, actx)
	return c
}

// GetActorContext returns additional actor.Context from application context
func (c *appContext) GetActorContext() actor.Context {
	rawVal, ok := c.Value(actorCtxKey).(actor.Context)
	if ok {
		return rawVal
	}
	return nil
}
