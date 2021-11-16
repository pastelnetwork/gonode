package storagechallenge

import (
	"time"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

type Config struct {
	common.Config                   `mapstructure:",squash" json:"-"`
	StorageChallengeExpiredDuration time.Duration `mapstructure:"storage_challenge_expired_duration" json:"storage_challenge_expired_duration"`
	NumberOfChallengeReplicas       int           `mapstructutr:"number_of_challenge_replicas" json:"number_of_challenge_replicas"`
}
