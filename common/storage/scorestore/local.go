package scorestore

import (
	"context"
)

// ScoreStorageInterface is interface for queries score sqlite store
type ScoreStorageInterface interface {
	CloseDB(ctx context.Context)

	ScoreAggregationTrackerQueries
	AggregateScScoreQueries
}
