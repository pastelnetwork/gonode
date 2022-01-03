package storagechallenge

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/pastel"
)

func Test_service_GenerateStorageChallenges(t *testing.T) {
	type fields struct {
		actor                         messaging.Actor
		domainActorID                 *actor.PID
		nodeID                        string
		pclient                       pastel.Client
		storageChallengeExpiredBlocks int32
		numberOfChallengeReplicas     int
		repository                    repository
		currentBlockCount             int32
	}
	type args struct {
		ctx                       context.Context
		challengesPerNodePerBlock int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				actor:                         tt.fields.actor,
				domainActorID:                 tt.fields.domainActorID,
				nodeID:                        tt.fields.nodeID,
				pclient:                       tt.fields.pclient,
				storageChallengeExpiredBlocks: tt.fields.storageChallengeExpiredBlocks,
				numberOfChallengeReplicas:     tt.fields.numberOfChallengeReplicas,
				repository:                    tt.fields.repository,
				currentBlockCount:             tt.fields.currentBlockCount,
			}
			if err := s.GenerateStorageChallenges(tt.args.ctx, tt.args.challengesPerNodePerBlock); (err != nil) != tt.wantErr {
				t.Errorf("service.GenerateStorageChallenges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
