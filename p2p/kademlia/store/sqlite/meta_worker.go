package sqlite

import (
	"context"
	"fmt"

	"os"
	"path"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/log"
)

var (
	commitLastAccessedInterval = 60 * time.Second
	migrationMetaDB            = "data001-migration-meta.sqlite3"
	accessUpdateBufferSize     = 100000
	commitInsertsInterval      = 90 * time.Second
	metaSyncBatchSize          = 10000

	updateChannel chan UpdateMessage
	insertChannel chan UpdateMessage
)

func init() {
	updateChannel = make(chan UpdateMessage, accessUpdateBufferSize)
	insertChannel = make(chan UpdateMessage, accessUpdateBufferSize)
}

type UpdateMessages []UpdateMessage

// AccessUpdate holds the key and the last accessed time.
type UpdateMessage struct {
	Key            string
	LastAccessTime time.Time
	Size           int
}

// MigrationMetaStore manages database operations.
type MigrationMetaStore struct {
	db           *sqlx.DB
	p2pDataStore *sqlx.DB

	updateTicker *time.Ticker
	insertTicker *time.Ticker

	updates sync.Map
	inserts sync.Map
}

// NewMigrationMetaStore initializes the MigrationMetaStore.
func NewMigrationMetaStore(ctx context.Context, dataDir string) (*MigrationMetaStore, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return nil, fmt.Errorf("mkdir %q: %w", dataDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot create data folder: %w", err)
	}

	dbFile := path.Join(dataDir, migrationMetaDB)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	p2pDataStore, err := connectP2PDataStore(dataDir)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting p2p store from meta-migration store")
	}

	handler := &MigrationMetaStore{
		db:           db,
		p2pDataStore: p2pDataStore,
		updateTicker: time.NewTicker(commitLastAccessedInterval),
		insertTicker: time.NewTicker(commitInsertsInterval),
	}

	if err := handler.migrateMeta(); err != nil {
		log.P2P().WithContext(ctx).Errorf("cannot create meta table in sqlite database: %s", err.Error())
	}

	if err := handler.migrateMetaMigration(); err != nil {
		log.P2P().WithContext(ctx).Errorf("cannot create meta-migration table in sqlite database: %s", err.Error())
	}

	if err := handler.migrateMigration(); err != nil {
		log.P2P().WithContext(ctx).Errorf("cannot create migration table in sqlite database: %s", err.Error())
	}

	if handler.isMetaSyncRequired() {
		err := handler.syncMetaWithData(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error syncing meta with p2p data")
		}
	}

	go handler.startLastAccessedUpdateWorker(ctx)
	go handler.startInsertWorker(ctx)

	return handler, nil
}

func (d *MigrationMetaStore) migrateMeta() error {
	query := `
	CREATE TABLE IF NOT EXISTS meta (
		key TEXT PRIMARY KEY,
		last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		access_count INTEGER DEFAULT 0,
		data_size INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	); `
	if _, err := d.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'del_keys': %w", err)
	}

	return nil
}

func (d *MigrationMetaStore) migrateMetaMigration() error {
	query := `
	CREATE TABLE IF NOT EXISTS meta_migration (
		key TEXT,
		migration_id INTEGER,
		score INTEGER,
		is_migrated BOOLEAN,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (key, migration_id)
	);`
	if _, err := d.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'del_keys': %w", err)
	}

	return nil
}

func (d *MigrationMetaStore) migrateMigration() error {
	query := `
	CREATE TABLE IF NOT EXISTS migration (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_data_size INTEGER,
		migration_started_at TIMESTAMP,
		migration_finished_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := d.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'del_keys': %w", err)
	}

	return nil
}

func (d *MigrationMetaStore) isMetaSyncRequired() bool {
	var exists int
	query := `SELECT EXISTS(SELECT 1 FROM meta LIMIT 1);`

	err := d.db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 0
}

func (d *MigrationMetaStore) syncMetaWithData(ctx context.Context) error {
	query := `SELECT key, data, updated_at FROM data LIMIT ? OFFSET ?`
	var offset int

	for {
		rows, err := d.p2pDataStore.Queryx(query, metaSyncBatchSize, offset)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error querying p2p data store")
			return err
		}

		var batchUpdates []UpdateMessage
		found := false
		for rows.Next() {
			found = true
			var r Record

			if err := rows.Scan(&r.Key, &r.Data, &r.UpdatedAt); err != nil {
				log.WithContext(ctx).WithError(err).Error("Error scanning row from p2p data store")
				continue
			}

			dataSize := len(r.Data)
			batchUpdates = append(batchUpdates, UpdateMessage{
				Key:            r.Key,
				LastAccessTime: r.UpdatedAt,
				Size:           dataSize,
			})
		}
		rows.Close()

		if !found {
			break
		}

		// Send batch for insertion using the existing channel-based mechanism.
		PostKeysInsert(batchUpdates)

		offset += len(batchUpdates) // Move the offset forward by the number of items processed.
	}

	return nil
}

func connectP2PDataStore(dataDir string) (*sqlx.DB, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return nil, fmt.Errorf("mkdir %q: %w", dataDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot create data folder: %w", err)
	}

	dbFile := path.Join(dataDir, dbName)
	dataDb, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}
	dataDb.SetMaxOpenConns(200)
	dataDb.SetMaxIdleConns(10)

	return dataDb, nil
}

// PostAccessUpdate sends access updates to be handled by the worker.
func PostAccessUpdate(updates []string) {
	for _, update := range updates {
		select {
		case updateChannel <- UpdateMessage{
			Key:            update,
			LastAccessTime: time.Now(),
		}:
			// Inserted
		default:
			log.WithContext(context.Background()).Error("updateChannel is full, dropping update")
		}

	}
}

// startWorker listens for updates and commits them periodically.
func (d *MigrationMetaStore) startLastAccessedUpdateWorker(ctx context.Context) {
	for {
		select {
		case update := <-updateChannel:
			d.updates.Store(update.Key, update.LastAccessTime)

		case <-d.updateTicker.C:
			d.commitLastAccessedUpdates(ctx)
		case <-ctx.Done():
			log.WithContext(ctx).Info("Shutting down last accessed update worker")
			d.commitLastAccessedUpdates(ctx) // Commit any remaining updates before shutdown
			return
		}
	}
}

func (d *MigrationMetaStore) commitLastAccessedUpdates(ctx context.Context) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error starting transaction (commitLastAccessedUpdates)")
		return
	}

	stmt, err := tx.Prepare("UPDATE meta SET last_access_time = ? WHERE key = ?")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error preparing statement (commitLastAccessedUpdates)")
		return
	}
	defer stmt.Close()

	keysToUpdate := make(map[string]time.Time)
	d.updates.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			return false
		}
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		_, err := stmt.Exec(v, k)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("key", key).Error("Error executing statement (commitLastAccessedUpdates)")
			return true // continue
		}
		keysToUpdate[k] = v

		return true // continue
	})

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.WithContext(ctx).WithError(err).Error("Error committing transaction (commitLastAccessedUpdates)")
		return
	}

	// Clear updated keys from map after successful commit
	for k := range keysToUpdate {
		d.updates.Delete(k)
	}

	log.WithContext(ctx).WithField("count", len(keysToUpdate)).Info("Committed last accessed updates")
}

func PostKeysInsert(updates []UpdateMessage) {
	for _, update := range updates {
		select {
		case insertChannel <- update:
			// Inserted
		default:
			log.WithContext(context.Background()).Error("insertChannel is full, dropping update")
		}
	}
}

// startInsertWorker listens for updates and commits them periodically.
func (d *MigrationMetaStore) startInsertWorker(ctx context.Context) {
	for {
		select {
		case update := <-insertChannel:
			d.inserts.Store(update.Key, update)
		case <-d.insertTicker.C:
			d.commitInserts(ctx)
		case <-ctx.Done():
			log.WithContext(ctx).Info("Shutting down insert meta keys worker")
			d.commitInserts(ctx)
			return
		}
	}
}

func (d *MigrationMetaStore) commitInserts(ctx context.Context) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error starting transaction (commitInserts)")
		return
	}

	// Prepare an INSERT OR REPLACE statement that handles new insertions or updates existing entries
	stmt, err := tx.Preparex("INSERT OR REPLACE INTO meta (key, last_accessed, access_count, data_size) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback() // Ensure to rollback in case of an error
		log.WithContext(ctx).WithError(err).Error("Error preparing statement (commitInserts)")
		return
	}
	defer stmt.Close()

	keysToUpdate := make(map[string]bool)
	d.inserts.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			return false
		}
		v, ok := value.(UpdateMessage)
		if !ok {
			return false
		}
		// Default values for access_count and data_size can be configured here
		accessCount := 1
		_, err := stmt.Exec(k, v.LastAccessTime, accessCount, v.Size)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("key", k).Error("Error executing statement (commitInserts)")
			return true // continue
		}
		keysToUpdate[k] = true

		return true // continue
	})

	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback transaction if commit fails
		log.WithContext(ctx).WithError(err).Error("Error committing transaction (commitInserts)")
		return
	}

	// Clear updated keys from map after successful commit
	for k := range keysToUpdate {
		d.inserts.Delete(k)
	}

	log.WithContext(ctx).WithField("count", len(keysToUpdate)).Info("Committed inserts")
}
