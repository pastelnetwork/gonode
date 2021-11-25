package context

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Context struct {
	context.Context
}

func (Context) GetActorContext() actor.Context {
	return nil
}
