package store

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

// DDStore represents Dupedetection store
type DDStore interface {
	GetLatestFingerprints(context.Context) (*domain.DDFingerprints, error)
	IfFingerprintExists(ctx context.Context, hash string) (bool, error)
	StoreFingerprint(context.Context, *domain.DDFingerprints) error
	GetFingerprintsCount(context.Context) (int64, error)
	CheckNonSeedRecord(ctx context.Context) (bool, error)
	IfFingerprintExistsByRegTxid(_ context.Context, txid string) (bool, error)
}

// ScoreStore is SN Score store
type ScoreStore interface {
	GetScoreByTxID(ctx context.Context, txid string) (*domain.SnScore, error)
	IncrementScore(ctx context.Context, score *domain.SnScore, increment int) (*domain.SnScore, error)
}

// PastelBlockStore is Pastel Block Store
type PastelBlockStore interface {
	StorePastelBlock(context.Context, domain.PastelBlock) error
	GetLatestPastelBlock(ctx context.Context) (domain.PastelBlock, error)
	GetPastelBlockByHash(ctx context.Context, hash string) (domain.PastelBlock, error)
	GetPastelBlockByHeight(ctx context.Context, height int32) (domain.PastelBlock, error)
	UpdatePastelBlock(ctx context.Context, block domain.PastelBlock) error
	FetchAllTxIDs() (map[string]bool, error)
	UpdateTxIDTimestamp(registrationTicketTxID string) error
}

// CollectionStore is collection store
type CollectionStore interface {
	IfCollectionExists(ctx context.Context, collectionTxID string) (bool, error)
	StoreCollection(_ context.Context, c domain.Collection) error
	GetCollection(ctx context.Context, collectionTxID string) (*domain.Collection, error)
	GetDoesNotImpactCollections(ctx context.Context, hash string) (domain.NonImpactedCollections, error)
	GetAllInProcessCollections(ctx context.Context) ([]*domain.Collection, error)
	FinalizeCollectionState(ctx context.Context, txid string) error
}

// SQLiteStore is sqlite implementation of DD store and Score store
type SQLiteStore struct {
	db *sqlx.DB
}

// NewSQLiteStore is new sqlite store constructor
func NewSQLiteStore(file string) (*SQLiteStore, error) {
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

	_, _ = db.Exec(createPbTableStatement)
	_, _ = db.Exec(createCollectionsTableStatement)
	_, _ = db.Exec(createFgTableStatement)
	_, _ = db.Exec(createDoesNotImpactCollectionsTableStatement)
	_, err = db.Exec(createScoreTableStatement)
	if err != nil {
		log.WithContext(context.Background()).WithError(err).Error("cannot create score table")
	}

	return &SQLiteStore{
		db: db,
	}, nil
}
