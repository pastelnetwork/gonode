package debug

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/storagechallenge"
	"github.com/pastelnetwork/gonode/pastel"
	"net/http"
	"time"
)

// processSCSummaryStats handles the core logic separately.
func (service *Service) processSCSummaryStats(ctx context.Context, pid string, signature string, from time.Time, _ *time.Time) (storagechallenge.SCSummaryStatsRes, error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("invalid pid/passphrase")
	}

	// Rate Limit Check for pid
	if !service.rateLimiter.CheckRateLimit(pid) {
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("rate limit exceeded, please try again later")
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.GetSCSummaryStats(from)
	if err != nil {
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("error querying metrics: %w", err)
	}

	return storagechallenge.SCSummaryStatsRes{
		SCSummaryStats: storagechallenge.SCSummaryStats{
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

// SCDetailedLogs is the HTTP handler that parses the request and writes the response.
func (service *Service) SCDetailedLogs(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	challengeID := request.URL.Query().Get("challenge_id")

	var storageChallengeDetailedLogsData []types.Message
	var err error
	if challengeID != "" {
		storageChallengeDetailedLogsData, err = service.GetSCDetailedLogsData(ctx, pid, passphrase, challengeID)
	} else {
		storageChallengeDetailedLogsData, err = service.GetNSCDetailedLogsData(ctx, pid, passphrase)
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

	responseWithJSON(writer, http.StatusOK, storageChallengeDetailedLogsData)
}

// GetSCDetailedLogsData encapsulates the core logic for storage-challenge log data
func (service *Service) GetSCDetailedLogsData(ctx context.Context, pid string, signature string, challengeID string) (storageChallengeMessageData []types.Message, err error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("invalid pid/passphrase")
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	if challengeID != "" {
		storageChallengeMessageData, err = store.GetMetricsDataByStorageChallengeID(ctx, challengeID)
		if err != nil {
			return storageChallengeMessageData, fmt.Errorf("error retrieving detailed logs: %w", err)
		}

		return storageChallengeMessageData, nil
	}

	return storageChallengeMessageData, nil
}

// GetNSCDetailedLogsData encapsulates the core logic for storage-challenge log data
func (service *Service) GetNSCDetailedLogsData(ctx context.Context, pid string, signature string) (storageChallengeMessageData []types.Message, err error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("invalid pid/passphrase")
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	mostRecentChallengeIDs, err := store.GetLastNSCMetrics()
	if err != nil {
		return storageChallengeMessageData, fmt.Errorf("error retrieving detailed logs: %w", err)
	}

	for _, mc := range mostRecentChallengeIDs {
		data, err := store.GetMetricsDataByStorageChallengeID(ctx, mc.ChallengeID)
		if err != nil {
			return storageChallengeMessageData, fmt.Errorf("error retrieving detailed logs: %w", err)
		}

		storageChallengeMessageData = append(storageChallengeMessageData, data...)
	}

	if len(storageChallengeMessageData) > 64 {
		return storageChallengeMessageData[:64], nil
	}

	return storageChallengeMessageData, nil
}

// SCSummaryStats is the HTTP handler that parses the request and writes the response.
func (service *Service) SCSummaryStats(writer http.ResponseWriter, request *http.Request) {
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

	storageChallengeSummaryStats, err := service.processSCSummaryStats(ctx, pid, passphrase, from, to)
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

	responseWithJSON(writer, http.StatusOK, storageChallengeSummaryStats)
}
