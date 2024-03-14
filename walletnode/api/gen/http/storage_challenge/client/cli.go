// Code generated by goa v3.15.0, DO NOT EDIT.
//
// StorageChallenge HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	storagechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge"
	goa "goa.design/goa/v3/pkg"
)

// BuildGetSummaryStatsPayload builds the payload for the StorageChallenge
// getSummaryStats endpoint from CLI flags.
func BuildGetSummaryStatsPayload(storageChallengeGetSummaryStatsFrom string, storageChallengeGetSummaryStatsTo string, storageChallengeGetSummaryStatsPid string, storageChallengeGetSummaryStatsKey string) (*storagechallenge.GetSummaryStatsPayload, error) {
	var err error
	var from *string
	{
		if storageChallengeGetSummaryStatsFrom != "" {
			from = &storageChallengeGetSummaryStatsFrom
			err = goa.MergeErrors(err, goa.ValidateFormat("from", *from, goa.FormatDateTime))
			if err != nil {
				return nil, err
			}
		}
	}
	var to *string
	{
		if storageChallengeGetSummaryStatsTo != "" {
			to = &storageChallengeGetSummaryStatsTo
			err = goa.MergeErrors(err, goa.ValidateFormat("to", *to, goa.FormatDateTime))
			if err != nil {
				return nil, err
			}
		}
	}
	var pid string
	{
		pid = storageChallengeGetSummaryStatsPid
	}
	var key string
	{
		key = storageChallengeGetSummaryStatsKey
	}
	v := &storagechallenge.GetSummaryStatsPayload{}
	v.From = from
	v.To = to
	v.Pid = pid
	v.Key = key

	return v, nil
}

// BuildGetDetailedLogsPayload builds the payload for the StorageChallenge
// getDetailedLogs endpoint from CLI flags.
func BuildGetDetailedLogsPayload(storageChallengeGetDetailedLogsPid string, storageChallengeGetDetailedLogsChallengeID string, storageChallengeGetDetailedLogsKey string) (*storagechallenge.GetDetailedLogsPayload, error) {
	var pid string
	{
		pid = storageChallengeGetDetailedLogsPid
	}
	var challengeID *string
	{
		if storageChallengeGetDetailedLogsChallengeID != "" {
			challengeID = &storageChallengeGetDetailedLogsChallengeID
		}
	}
	var key string
	{
		key = storageChallengeGetDetailedLogsKey
	}
	v := &storagechallenge.GetDetailedLogsPayload{}
	v.Pid = pid
	v.ChallengeID = challengeID
	v.Key = key

	return v, nil
}