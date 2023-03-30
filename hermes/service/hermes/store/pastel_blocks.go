package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
)

const (
	createPbTableStatement = `CREATE TABLE IF NOT EXISTS pastel_blocks_table (block_hash text PRIMARY KEY, block_height int, datetime_block_added text)`
)

type pastelBlock struct {
	BlockHash          string `db:"block_hash"`
	BlockHeight        int32  `db:"block_height"`
	DatetimeBlockAdded string `db:"datetime_block_added"`
}

func (p pastelBlock) toDomain() domain.PastelBlock {
	return domain.PastelBlock{
		BlockHash:          p.BlockHash,
		BlockHeight:        p.BlockHeight,
		DatetimeBlockAdded: p.DatetimeBlockAdded,
	}
}

// GetLatestPastelBlock gets the latest record from pastel_blocks table
func (s *SQLiteStore) GetLatestPastelBlock(ctx context.Context) (domain.PastelBlock, error) {
	pb := pastelBlock{}

	getLatestPastelBlockQuery := `SELECT * from pastel_blocks_table order by datetime_block_added DESC LIMIT 1;`
	err := s.db.GetContext(ctx, &pb, getLatestPastelBlockQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			return pb.toDomain(), nil
		}

		return pb.toDomain(), fmt.Errorf("failed to get record by key %w", err)
	}

	return pb.toDomain(), nil
}

// StorePastelBlock store the block hash and height in pastel_blocks table
func (s *SQLiteStore) StorePastelBlock(ctx context.Context, pb domain.PastelBlock) error {
	pastelBlockQuery := `INSERT into pastel_blocks_table(block_hash, block_height, datetime_block_added) VALUES(?,?,?)`

	_, err := s.db.ExecContext(ctx, pastelBlockQuery, pb.BlockHash, pb.BlockHeight, pb.DatetimeBlockAdded)
	if err != nil {
		return errors.Errorf("unable to insert into pastel blocks table")
	}

	return nil
}