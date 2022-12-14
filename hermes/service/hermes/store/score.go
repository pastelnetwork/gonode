package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //go-sqlite3

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
)

const (
	createScoreTableStatement = `CREATE TABLE IF NOT EXISTS sn_scores (txid text PRIMARY KEY, ext_key text, ip text, score integer)`
	getScoreFromTxidStatement = `SELECT * FROM sn_scores WHERE txid = ?`
	insertSnScoreStatement    = `INSERT INTO sn_scores(txid, ext_key, ip, score) VALUES(?,?,?,?)`
	updateSnScoreStatement    = `UPDATE sn_scores SET score = ? WHERE txid = ?`
)

type snScore struct {
	TxID      string `db:"txid,omitempty"`
	ExtKey    string `db:"ext_key,omitempty"`
	IPAddress string `db:"ip,omitempty"`
	Score     int    `db:"score,omitempty"`
}

func (r *snScore) toDomain() *domain.SnScore {

	return &domain.SnScore{
		TxID:      r.TxID,
		PastelID:  r.ExtKey,
		IPAddress: r.IPAddress,
		Score:     r.Score,
	}
}

// GetScoreByTxID gets score by txid
func (s *SQLiteStore) GetScoreByTxID(_ context.Context, txid string) (*domain.SnScore, error) {
	r := snScore{}
	err := s.db.Get(&r, getScoreFromTxidStatement, txid)
	if err != nil {
		return nil, fmt.Errorf("failed to get sn_score record by key %w : key: %s", err, txid)
	}

	return r.toDomain(), nil
}

// IncrementScore increments score and if not exists, it creates a record
func (s *SQLiteStore) IncrementScore(ctx context.Context, score *domain.SnScore, increment int) (*domain.SnScore, error) {
	r := snScore{}
	err := s.db.Get(&r, getScoreFromTxidStatement, score.TxID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.db.Exec(insertSnScoreStatement, score.TxID, score.PastelID, score.IPAddress, score.Score)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("txid", score.TxID).Error("Failed to insert sn_score record")
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("failed to get sn_score record by key %w : key: %s", err, score.TxID)
		}
	}

	r.Score = r.Score + increment
	_, err = s.db.Exec(updateSnScoreStatement, r.Score, r.TxID)
	if err != nil {
		log.WithContext(ctx).WithField("txid", r.TxID).WithError(err).Error("Failed to update sn_score record")
		return nil, err
	}

	return r.toDomain(), nil
}
