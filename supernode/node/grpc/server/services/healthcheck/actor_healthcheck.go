package healthcheck

import (
	"context"
	"encoding/json"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
)

// Receive represents actor mocel receive action
func (service *HealthCheck) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *pb.PingRequest:
		stats, err := service.Stats(context.Background())
		if err != nil {
			ctx.Respond(errors.Errorf("Stats(): %w", err))
			return
		}

		jsonData, err := json.Marshal(stats)
		if err != nil {
			ctx.Respond(errors.Errorf("Marshal(): %w", err))
			return
		}

		// echos received message
		ctx.Respond(&pb.PingReply{Reply: string(jsonData)})
	}
}
