// Code generated by goa v3.15.0, DO NOT EDIT.
//
// Score HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	score "github.com/pastelnetwork/gonode/walletnode/api/gen/score"
)

// BuildGetAggregatedChallengesScoresPayload builds the payload for the Score
// getAggregatedChallengesScores endpoint from CLI flags.
func BuildGetAggregatedChallengesScoresPayload(scoreGetAggregatedChallengesScoresPid string, scoreGetAggregatedChallengesScoresKey string) (*score.GetAggregatedChallengesScoresPayload, error) {
	var pid string
	{
		pid = scoreGetAggregatedChallengesScoresPid
	}
	var key string
	{
		key = scoreGetAggregatedChallengesScoresKey
	}
	v := &score.GetAggregatedChallengesScoresPayload{}
	v.Pid = pid
	v.Key = key

	return v, nil
}