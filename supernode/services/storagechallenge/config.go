package storagechallenge

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	//Only one block should pass in the time between storage challenge and verification. Currently block time is 2.5 minutes so this should be enough.
	defaultStorageChallengeExpiredBlocks = 1
	//Number of challenge replicas defaults to 3
	defaultNumberOfChallengeReplicas = 3
	//Number of verifying nodes determines how many potential verifiers we test in verify_storage_challenge
	defaultNumberOfVerifyingNodes = 10
)

// Config storage challenge config
type Config struct {
	common.Config                 `mapstructure:",squash" json:"-"`
	StorageChallengeExpiredBlocks int32 `mapstructure:"storage_challenge_expired_blocks" json:"storage_challenge_expired_duration"`
	NumberOfChallengeReplicas     int   `mapstructure:"number_of_challenge_replicas" json:"number_of_challenge_replicas"`
	NumberOfVerifyingNodes        int   `mapstructure:"number_of_verifying_nodes" json:"number_of_verifying_nodes"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                        *common.NewConfig(),
		StorageChallengeExpiredBlocks: defaultStorageChallengeExpiredBlocks,
		NumberOfChallengeReplicas:     defaultNumberOfChallengeReplicas,
		NumberOfVerifyingNodes:        defaultNumberOfVerifyingNodes,
	}
}
