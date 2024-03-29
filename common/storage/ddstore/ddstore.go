package ddstore

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	getFingerprintFromHashStatement = `SELECT sha256_hash_of_art_image_file FROM image_hash_to_image_fingerprint_table WHERE sha256_hash_of_art_image_file = ?`
)

// DDStore represents Dupedetection store
type DDStore interface {
	// GetDDDataHash returns hash of dd data
	GetDDDataHash(ctx context.Context) (hash string, err error)
	// IfFingerprintExists checks if fg exists against the hash
	IfFingerprintExists(_ context.Context, hash string) (bool, error)
}

// SQLiteDDStore is sqlite implementation of DD store and Score store
type SQLiteDDStore struct {
	db *sqlx.DB
}

// NewSQLiteDDStore is new sqlite store constructor
func NewSQLiteDDStore(file string) (*SQLiteDDStore, error) {
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		tmpfile, err := ioutil.TempFile("", "registered_image_fingerprints_db.sqlite")
		if err != nil {
			panic(err.Error())
		}
		file = tmpfile.Name()
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.Errorf("database dd service not found: %w", err)
	}

	db, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("cannot open dd-service database: %w", err)
	}

	return &SQLiteDDStore{
		db: db,
	}, nil
}

type fingerprints struct {
	Sha256HashOfArtImageFile string `db:"sha256_hash_of_art_image_file,omitempty"`
}
type collections struct {
	CollectionTicketTXID string `db:"collection_ticket_txid"`
}

// IfFingerprintExists checks if fg exists against the hash
func (s *SQLiteDDStore) IfFingerprintExists(_ context.Context, hash string) (bool, error) {
	r := fingerprints{}
	err := s.db.Get(&r, getFingerprintFromHashStatement, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("failed to get record by key %w : key: %s", err, hash)
	}

	if r.Sha256HashOfArtImageFile == "" {
		return false, nil
	}

	return true, nil
}

// GetDDDataHash returns hash of dd data
func (s *SQLiteDDStore) GetDDDataHash(ctx context.Context) (hash string, err error) {
	r := []fingerprints{}
	err = s.db.Select(&r, "SELECT sha256_hash_of_art_image_file FROM image_hash_to_image_fingerprint_table where txid_timestamp > 0 order by sha256_hash_of_art_image_file asc")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get image_hash_to_image_fingerprint_table")
	}

	c := []collections{}
	err = s.db.Select(&c, "SELECT collection_ticket_txid FROM collections_table")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get collections_table, ignore table")
	}

	var sb strings.Builder

	for i := 0; i < len(r); i++ {
		sb.WriteString(r[i].Sha256HashOfArtImageFile)
	}

	for i := 0; i < len(c); i++ {
		sb.WriteString(c[i].CollectionTicketTXID)
	}

	hash = utils.GetHashFromString(sb.String())
	log.WithContext(ctx).WithField("hash", hash).Debug("dd data hash returned")

	return hash, nil
}

// Close closes the database
func (s *SQLiteDDStore) Close() error {
	return s.db.Close()
}
