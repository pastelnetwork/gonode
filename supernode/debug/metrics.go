package debug

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// processSHChallenge encapsulates the core logic for self-healing challenge processing.
func (service *Service) processSHChallenge(ctx context.Context, pid string, passphrase string, challengeID string, count int) (reports types.SelfHealingChallengeReports, err error) {
	_, err = service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("invalid pid/passphrase: %w", err)
	}
	reports = make(types.SelfHealingChallengeReports)

	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	if challengeID != "" {
		reports, err = store.GetSHChallengeReport(ctx, challengeID)
		if err != nil {
			return reports, fmt.Errorf("error retrieving report: %w", err)
		}

		return reports, nil

	}

	reports, err = store.GetLastNSHChallenges(ctx, count)
	if err != nil {
		return reports, fmt.Errorf("error retrieving report: %w", err)
	}

	return reports, nil

}

// shChallenge is the HTTP handler for self-healing challenges.
func (service *Service) shChallenge(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	challengeID := request.URL.Query().Get("challenge_id")
	// Parse the optional count parameter
	countStr := request.URL.Query().Get("count")
	count := 10 // Default value of 10 if count is not provided or invalid
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid count parameter"})
			return
		}
	}

	result, err := service.processSHChallenge(ctx, pid, passphrase, challengeID, count)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "invalid pid/passphrase":
			statusCode = http.StatusUnauthorized
		case "error opening DB", "error retrieving metrics":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusNotFound // For "no metrics found for challenge_id"
		}
		responseWithJSON(writer, statusCode, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, result)
}

// processSHTrigger encapsulates the core logic for self-healing trigger processing.
func (service *Service) processSHTrigger(ctx context.Context, pid string, passphrase string, triggerID string, count int) ([]types.SelfHealingGenerationMetric, error) {
	_, err := service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.GetSelfHealingGenerationMetrics(time.Now().AddDate(20, 0, 0))
	if err != nil {
		return nil, fmt.Errorf("error retrieving metrics: %w", err)
	}

	if triggerID == "" {
		if len(metrics) > count {
			return metrics[:count], nil
		}

		return metrics, nil
	}

	for _, metric := range metrics {
		if metric.TriggerID == triggerID {
			return []types.SelfHealingGenerationMetric{metric}, nil
		}
	}

	return nil, fmt.Errorf("no metrics found for trigger_id %s", triggerID)
}

// shTrigger is the HTTP handler for self-healing triggers.
func (service *Service) shTrigger(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	triggerID := request.URL.Query().Get("trigger_id")

	// Parse the optional count parameter
	countStr := request.URL.Query().Get("count")
	count := 10 // Default value of 10 if count is not provided or invalid
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid count parameter"})
			return
		}
	}

	result, err := service.processSHTrigger(ctx, pid, passphrase, triggerID, count)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "invalid pid/passphrase":
			statusCode = http.StatusUnauthorized
		case "error opening DB", "error retrieving metrics":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusNotFound // For "no metrics found for trigger_id"
		}
		responseWithJSON(writer, statusCode, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, result)
}

// processMetrics handles the core logic separately.
func (service *Service) processMetrics(ctx context.Context, pid string, passphrase string, from time.Time, to *time.Time) (interface{}, error) {
	_, err := service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	// Rate Limit Check for pid
	if !service.rateLimiter.CheckRateLimit(pid) {
		return nil, fmt.Errorf("rate limit exceeded, please try again later")
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.QueryMetrics(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("error querying metrics: %w", err)
	}

	return metrics, nil
}

// metrics is the HTTP handler that parses the request and writes the response.
func (service *Service) metrics(writer http.ResponseWriter, request *http.Request) {
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

	metrics, err := service.processMetrics(ctx, pid, passphrase, from, to)
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

	responseWithJSON(writer, http.StatusOK, metrics)
}
