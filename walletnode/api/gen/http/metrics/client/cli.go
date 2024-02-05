// Code generated by goa v3.14.0, DO NOT EDIT.
//
// metrics HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"fmt"
	"strconv"

	metrics "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	goa "goa.design/goa/v3/pkg"
)

// BuildGetChallengeReportsPayload builds the payload for the metrics
// getChallengeReports endpoint from CLI flags.
func BuildGetChallengeReportsPayload(metricsGetChallengeReportsPid string, metricsGetChallengeReportsChallengeID string, metricsGetChallengeReportsCount string, metricsGetChallengeReportsKey string) (*metrics.GetChallengeReportsPayload, error) {
	var err error
	var pid string
	{
		pid = metricsGetChallengeReportsPid
	}
	var challengeID *string
	{
		if metricsGetChallengeReportsChallengeID != "" {
			challengeID = &metricsGetChallengeReportsChallengeID
		}
	}
	var count *int
	{
		if metricsGetChallengeReportsCount != "" {
			var v int64
			v, err = strconv.ParseInt(metricsGetChallengeReportsCount, 10, strconv.IntSize)
			val := int(v)
			count = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for count, must be INT")
			}
		}
	}
	var key string
	{
		key = metricsGetChallengeReportsKey
	}
	v := &metrics.GetChallengeReportsPayload{}
	v.Pid = pid
	v.ChallengeID = challengeID
	v.Count = count
	v.Key = key

	return v, nil
}

// BuildGetMetricsPayload builds the payload for the metrics getMetrics
// endpoint from CLI flags.
func BuildGetMetricsPayload(metricsGetMetricsFrom string, metricsGetMetricsTo string, metricsGetMetricsPid string, metricsGetMetricsKey string) (*metrics.GetMetricsPayload, error) {
	var err error
	var from *string
	{
		if metricsGetMetricsFrom != "" {
			from = &metricsGetMetricsFrom
			err = goa.MergeErrors(err, goa.ValidateFormat("from", *from, goa.FormatDateTime))
			if err != nil {
				return nil, err
			}
		}
	}
	var to *string
	{
		if metricsGetMetricsTo != "" {
			to = &metricsGetMetricsTo
			err = goa.MergeErrors(err, goa.ValidateFormat("to", *to, goa.FormatDateTime))
			if err != nil {
				return nil, err
			}
		}
	}
	var pid string
	{
		pid = metricsGetMetricsPid
	}
	var key string
	{
		key = metricsGetMetricsKey
	}
	v := &metrics.GetMetricsPayload{}
	v.From = from
	v.To = to
	v.Pid = pid
	v.Key = key

	return v, nil
}
