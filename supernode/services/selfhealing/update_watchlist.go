package selfhealing

import (
	"context"
	"database/sql"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

// UpdateWatchlist fetch and update the nodes on watchlist
func (task *SHTask) UpdateWatchlist(ctx context.Context) error {
	log.WithContext(ctx).Infoln("Update Watchlist worker invoked")

	pingInfos, err := task.GetAllPingInfo(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping infos")
	}

	for _, info := range pingInfos {

		if task.shouldUpdateWatchlistField(info) {
			info.IsOnWatchlist = true

			if err := task.UpsertPingHistory(ctx, info); err != nil {
				log.WithContext(ctx).
					WithField("supernode_id", info.SupernodeID).
					WithError(err).
					Error("error upserting ping history")
			}

		}
	}

	return nil
}

// GetAllPingInfo get all the ping info from db
func (task *SHTask) GetAllPingInfo(ctx context.Context) (types.PingInfos, error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return nil, err
	}

	var infos types.PingInfos

	if store != nil {
		defer store.CloseHistoryDB(ctx)

		infos, err = store.GetAllPingInfos()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			log.WithContext(ctx).
				Error("error retrieving ping histories")

			return nil, err
		}
	}

	return infos, nil
}

// UpsertPingHistory upsert the ping info to db
func (task *SHTask) UpsertPingHistory(ctx context.Context, info types.PingInfo) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}

	if store != nil {
		defer store.CloseHistoryDB(ctx)

		err = store.UpsertPingHistory(info)
		if err != nil {
			log.WithContext(ctx).
				WithField("supernode_id", info.SupernodeID).
				Error("error upserting ping history")

			return err
		}

	}

	return nil
}

func (task *SHTask) shouldUpdateWatchlistField(info types.PingInfo) bool {
	twentyMinutesAgo := time.Now().UTC().Add(-20 * time.Minute)

	// Check if the timestamp is before 20 minutes ago
	if info.LastSeen.Time.Before(twentyMinutesAgo) {
		return true
	}

	return false

}
