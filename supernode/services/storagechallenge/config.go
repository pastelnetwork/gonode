package storagechallenge

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// Config storage challenge config
type Config struct {
	common.Config                 `mapstructure:",squash" json:"-"`
	StorageChallengeExpiredBlocks int32 `mapstructure:"storage_challenge_expired_blocks" json:"storage_challenge_expired_duration"`
	NumberOfChallengeReplicas     int   `mapstructutr:"number_of_challenge_replicas" json:"number_of_challenge_replicas"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
}
