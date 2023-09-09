package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/hermes/domain"
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

// GetPastelBlockByHash gets the pastel-block by hash
func (s *SQLiteStore) GetPastelBlockByHash(ctx context.Context, hash string) (domain.PastelBlock, error) {
	pb := pastelBlock{}

	getPastelBlockByHash := `SELECT * from pastel_blocks_table where block_hash=?;`
	err := s.db.GetContext(ctx, &pb, getPastelBlockByHash, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return pb.toDomain(), nil
		}

		return pb.toDomain(), fmt.Errorf("failed to get record by key %w", err)
	}

	return pb.toDomain(), nil
}

// GetPastelBlockByHeight gets the pastel-block by height
func (s *SQLiteStore) GetPastelBlockByHeight(ctx context.Context, height int32) (domain.PastelBlock, error) {
	pb := pastelBlock{}

	getPastelBlockByHeight := `SELECT * from pastel_blocks_table where block_height=?;`
	err := s.db.GetContext(ctx, &pb, getPastelBlockByHeight, height)
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Errorf("unable to begin transaction")
	}

	pastelBlockQuery := `INSERT OR IGNORE into pastel_blocks_table(block_hash, block_height, datetime_block_added) VALUES(?,?,?)`
	_, err = tx.ExecContext(ctx, pastelBlockQuery, pb.BlockHash, pb.BlockHeight, pb.DatetimeBlockAdded)
	if err != nil {
		tx.Rollback()
		return errors.Errorf("unable to insert into pastel blocks table: height: %d - hash: %s - err: %w", pb.BlockHeight, pb.BlockHash, err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Errorf("unable to commit txn: height: %d - hash: %s - err: %w", pb.BlockHeight, pb.BlockHash, err)
	}

	return nil
}

// UpdatePastelBlock updates the pastel-block by height
func (s *SQLiteStore) UpdatePastelBlock(ctx context.Context, pb domain.PastelBlock) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Errorf("unable to begin transaction : %w", err)
	}

	updatePastelBlockQuery := `update pastel_blocks_table set block_hash=? where block_height =?`
	_, err = tx.ExecContext(ctx, updatePastelBlockQuery, pb.BlockHash, pb.BlockHeight)
	if err != nil {
		tx.Rollback()
		return errors.Errorf("unable to update pastel blocks table by height: %w", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Errorf("unable to commit transaction : %w", err)
	}

	return nil
}

func (s *SQLiteStore) FetchAllTxIDs() (map[string]bool, error) {
	txIDs := make(map[string]bool)

	rows, err := s.db.Query("SELECT registration_ticket_txid FROM image_hash_to_image_fingerprint_table WHERE registration_ticket_txid IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txID string
	for rows.Next() {
		if err := rows.Scan(&txID); err != nil {
			return nil, err
		}
		txIDs[txID] = true
	}

	return txIDs, rows.Err()
}

func (s *SQLiteStore) UpdateTxIDTimestamp(registrationTicketTxID string) error {
	query := "UPDATE image_hash_to_image_fingerprint_table SET txid_timestamp = 0 WHERE registration_ticket_txid = ?"

	_, err := s.db.Exec(query, registrationTicketTxID)
	return err
}
