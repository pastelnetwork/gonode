package storagechallenge

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

func TestTaskGenerateStorageChallenges(t *testing.T) {
	type fields struct {
		SuperNodeTask           *common.SuperNodeTask
		StorageChallengeService *StorageChallengeService
		storage                 *common.StorageHandler
	}
	type args struct {
		ctx context.Context
	}
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"success": {
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := StorageChallengeTask{
				SuperNodeTask:           tt.fields.SuperNodeTask,
				StorageChallengeService: tt.fields.StorageChallengeService,
				storage:                 tt.fields.storage,
			}
			if err := task.GenerateStorageChallenges(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("StorageChallengeTask.GenerateStorageChallenges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
