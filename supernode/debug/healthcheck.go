package debug

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
	healthCheckChallenge "github.com/pastelnetwork/gonode/common/utils/healthcheckchallenge"
	"github.com/pastelnetwork/gonode/pastel"
)

// processHCSummaryStats handles the core logic separately.
func (service *Service) processHCSummaryStats(ctx context.Context, pid string, signature string, from time.Time, _ *time.Time) (healthCheckChallenge.HCSummaryStatsRes, error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return healthCheckChallenge.HCSummaryStatsRes{}, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return healthCheckChallenge.HCSummaryStatsRes{}, fmt.Errorf("invalid pid/passphrase")
	}

	// Rate Limit Check for pid
	if !service.rateLimiter.CheckRateLimit(pid) {
		return healthCheckChallenge.HCSummaryStatsRes{}, fmt.Errorf("rate limit exceeded, please try again later")
	}

	store, err := queries.OpenHistoryDB()
	if err != nil {
		return healthCheckChallenge.HCSummaryStatsRes{}, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.GetHCSummaryStats(from)
	if err != nil {
		return healthCheckChallenge.HCSummaryStatsRes{}, fmt.Errorf("error querying metrics: %w", err)
	}

	return healthCheckChallenge.HCSummaryStatsRes{
		HCSummaryStats: healthCheckChallenge.HCSummaryStats{
			TotalChallenges:                      metrics.TotalChallenges,
			TotalChallengesProcessed:             metrics.TotalChallengesProcessed,
			TotalChallengesEvaluatedByChallenger: metrics.TotalChallengesEvaluatedByChallenger,
			TotalChallengesVerified:              metrics.TotalChallengesVerified,
			SlowResponsesObservedByObservers:     metrics.SlowResponsesObservedByObservers,
			InvalidSignaturesObservedByObservers: metrics.InvalidSignaturesObservedByObservers,
			InvalidEvaluationObservedByObservers: metrics.InvalidEvaluationObservedByObservers,
		},
	}, nil
}

// HCSummaryStats is the HTTP handler that parses the request and writes the response.
func (service *Service) HCSummaryStats(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	fromStr := request.URL.Query().Get("start")
	var from time.Time
	if fromStr == "" {
		from = time.Now().UTC().Add(-time.Hour * 24 * 7) // Default to 1 week ago
	} else {
		var err error
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid start time format"})
			return
		}
	}

	toStr := request.URL.Query().Get("end")
	var to *time.Time
	if toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid end time format"})
			return
		}
		to = &parsedTo
	}

	healthCheckChallengeSummaryStats, err := service.processHCSummaryStats(ctx, pid, passphrase, from, to)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "invalid pid/passphrase":
			statusCode = http.StatusUnauthorized
		case "rate limit exceeded, please try again later":
			statusCode = http.StatusTooManyRequests
		case "error opening DB", "error querying metrics":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		responseWithJSON(writer, statusCode, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, healthCheckChallengeSummaryStats)
}

// HCDetailedLogs is the HTTP handler that parses the request and writes the response.
func (service *Service) HCDetailedLogs(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	challengeID := request.URL.Query().Get("challenge_id")

	var HCDetailedLogsData []types.HealthCheckMessage
	var err error
	if challengeID != "" {
		HCDetailedLogsData, err = service.GetHCDetailedLogsData(ctx, pid, passphrase, challengeID)
	} else {
		HCDetailedLogsData, err = service.GetNHCDetailedLogsData(ctx, pid, passphrase)
	}

	if err != nil {
		var statusCode int
		switch err.Error() {
		case "invalid pid/passphrase":
			statusCode = http.StatusUnauthorized
		case "error opening DB", "error querying metrics":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		responseWithJSON(writer, statusCode, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, HCDetailedLogsData)

}

// GetHCDetailedLogsData encapsulates the core logic for storage-challenge log data
func (service *Service) GetHCDetailedLogsData(ctx context.Context, pid string, signature string, challengeID string) (healthCheckChallengeMessage []types.HealthCheckMessage, err error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("invalid pid/passphrase")
	}

	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	if challengeID != "" {
		healthCheckChallengeMessage, err = store.GetMetricsDataByHealthCheckChallengeID(ctx, challengeID)
		if err != nil {
			return healthCheckChallengeMessage, fmt.Errorf("error retrieving detailed logs: %w", err)
		}

		return healthCheckChallengeMessage, nil
	}

	return healthCheckChallengeMessage, nil
}

func (service *Service) GetNHCDetailedLogsData(ctx context.Context, pid string, signature string) (healthCheckMessageData []types.HealthCheckMessage, err error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("invalid pid/passphrase")
	}

	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	mostRecentChallengeIDs, err := store.GetLastNHCMetrics()
	if err != nil {
		return healthCheckMessageData, fmt.Errorf("error retrieving detailed logs: %w", err)
	}

	for _, mc := range mostRecentChallengeIDs {
		data, err := store.GetMetricsDataByHealthCheckChallengeID(ctx, mc.ChallengeID)
		if err != nil {
			return healthCheckMessageData, fmt.Errorf("error retrieving detailed logs: %w", err)
		}

		healthCheckMessageData = append(healthCheckMessageData, data...)
	}

	if len(healthCheckMessageData) > 64 {
		return healthCheckMessageData[:64], nil
	}

	return healthCheckMessageData, nil
}
