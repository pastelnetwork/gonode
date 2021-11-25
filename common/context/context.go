package context

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// Context application context wrapping base context
type Context struct {
	context.Context
}

// GetActorContext returns additional actor.Context from application context
func (Context) GetActorContext() actor.Context {
	return nil
}
