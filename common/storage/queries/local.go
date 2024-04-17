package queries

import (
	"context"
)

// LocalStoreInterface is interface for queries sqlite store
type LocalStoreInterface interface {
	CloseHistoryDB(ctx context.Context)

	TaskHistoryQueries
	SelfHealingQueries
	StorageChallengeQueries
	PingHistoryQueries
	HealthCheckChallengeQueries
}
