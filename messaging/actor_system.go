package messaging

import "github.com/AsynkronIT/protoactor-go/actor"

var actorSystem *actor.ActorSystem

// GetActorSystem get root actor system
func GetActorSystem() *actor.ActorSystem {
	if actorSystem == nil {
		actorSystem = actor.NewActorSystem()
	}
	return actorSystem
}
