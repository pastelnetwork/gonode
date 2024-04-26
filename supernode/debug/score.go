package debug

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/scorestore"
	"github.com/pastelnetwork/gonode/common/types"
	"net/http"

	"github.com/pastelnetwork/gonode/pastel"
)

// GetChallengesScore returns the challenges score
func (service *Service) GetChallengesScore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	passphrase := request.Header.Get("Authorization")
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	scores, err := service.getScores(ctx, pid, passphrase)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "invalid pid/passphrase":
			statusCode = http.StatusUnauthorized
		case "rate limit exceeded, please try again later":
			statusCode = http.StatusTooManyRequests
		case "error opening DB", "error querying scores":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		responseWithJSON(writer, statusCode, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, scores)
}

// getScores fetches the scores from the DB and return it
func (service *Service) getScores(ctx context.Context, pid string, signature string) (interface{}, error) {
	ok, err := service.scService.PastelClient.Verify(ctx, []byte(pid), signature, pid, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pid/passphrase: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("invalid pid/passphrase")
	}

	// Rate Limit Check for pid
	if !service.rateLimiter.CheckRateLimit(pid) {
		return nil, fmt.Errorf("rate limit exceeded, please try again later")
	}

	store, err := scorestore.OpenScoringDb()
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	defer store.CloseDB(ctx)

	supernodes, err := service.scService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return nil, err
	}

	var aggregatedScores types.AggregatedScoreList
	for _, node := range supernodes {
		var score *types.AggregatedScore

		score, err = store.GetAggregatedScore(node.ExtKey)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error querying challenges score: %w", err)
		}

		if score == nil {
			score = &types.AggregatedScore{
				NodeID:                    node.ExtKey,
				IPAddress:                 node.IPAddress,
				StorageChallengeScore:     50,
				HealthCheckChallengeScore: 50,
			}
		}

		aggregatedScores = append(aggregatedScores, *score)
	}

	return aggregatedScores, nil
}
