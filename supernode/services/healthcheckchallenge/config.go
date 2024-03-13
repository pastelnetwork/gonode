package healthcheckchallenge

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	//Only one block should pass in the time between healthcheck challenge and verification. Currently block time is 2.5 minutes so this should be enough.
	defaultHealthcheckChallengeExpiredBlocks = 1
	//Number of challenge replicas defaults to 6
	defaultNumberOfChallengeReplicas = 6
	//Number of verifying nodes determines how many potential verifiers we test in verify_healthcheck_challenge
	defaultNumberOfVerifyingNodes = 10
	// SuccessfulEvaluationThreshold is the affirmation threshold required for broadcasting
	SuccessfulEvaluationThreshold = 3
)

// Config health check challenge config
type Config struct {
	common.Config                     `mapstructure:",squash" json:"-"`
	HealthCheckChallengeExpiredBlocks int32 `mapstructure:"healthcheck_challenge_expired_blocks" json:"healthcheck_challenge_expired_duration"`
	NumberOfChallengeReplicas         int   `mapstructure:"number_of_challenge_replicas" json:"number_of_challenge_replicas"`
	NumberOfVerifyingNodes            int   `mapstructure:"number_of_verifying_nodes" json:"number_of_verifying_nodes"`

	// raptorq service
	RaptorQServiceAddress         string `mapstructure:"-" json:"-"`
	RqFilesDir                    string
	SuccessfulEvaluationThreshold int
	IsTestConfig                  bool
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                            *common.NewConfig(),
		HealthCheckChallengeExpiredBlocks: defaultHealthcheckChallengeExpiredBlocks,
		NumberOfChallengeReplicas:         defaultNumberOfChallengeReplicas,
		NumberOfVerifyingNodes:            defaultNumberOfVerifyingNodes,
		SuccessfulEvaluationThreshold:     SuccessfulEvaluationThreshold,
	}
}
