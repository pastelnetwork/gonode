package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"os"
	"path"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/cloud.go"
)

var (
	commitLastAccessedInterval = 60 * time.Second
	migrationExecutionTicker   = 12 * time.Hour
	migrationMetaDB            = "data001-migration-meta.sqlite3"
	accessUpdateBufferSize     = 100000
	commitInsertsInterval      = 90 * time.Second
	metaSyncBatchSize          = 10000
	lowSpaceThresholdGB        = 50 // in GB
	minKeysToMigrate           = 100

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
	cloud        cloud.Storage

	updateTicker             *time.Ticker
	insertTicker             *time.Ticker
	migrationExecutionTicker *time.Ticker

	updates sync.Map
	inserts sync.Map
}

// NewMigrationMetaStore initializes the MigrationMetaStore.
func NewMigrationMetaStore(ctx context.Context, dataDir string, cloud cloud.Storage) (*MigrationMetaStore, error) {
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
		db:                       db,
		p2pDataStore:             p2pDataStore,
		updateTicker:             time.NewTicker(commitLastAccessedInterval),
		insertTicker:             time.NewTicker(commitInsertsInterval),
		migrationExecutionTicker: time.NewTicker(migrationExecutionTicker),
		cloud:                    cloud,
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
	go handler.startMigrationExecutionWorker(ctx)

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
	query := `SELECT key, data, updatedAt FROM data LIMIT ? OFFSET ?`
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

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO meta (key, last_accessed) VALUES (?, ?)")
	if err != nil {
		tx.Rollback() // Roll back the transaction on error
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

// startMigrationExecutionWorker starts the worker that executes a migration
func (d *MigrationMetaStore) startMigrationExecutionWorker(ctx context.Context) {
	for {
		select {
		case <-d.migrationExecutionTicker.C:
			d.checkAndExecuteMigration(ctx)
		case <-ctx.Done():
			log.WithContext(ctx).Info("Shutting down data migration worker")
			return
		}
	}
}

func (d *MigrationMetaStore) checkAndExecuteMigration(ctx context.Context) {
	// Check the available disk space
	isLow, err := utils.CheckDiskSpace(lowSpaceThresholdGB)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("migration worker: check disk space failed")
	}

	if !isLow {
		// Disk space is sufficient, stop migration
		return
	}

	// Step 1: Fetch pending migrations
	var migrations []struct {
		ID int `db:"id"`
	}

	err = d.db.Select(&migrations, `SELECT id FROM migration WHERE migration_started_at IS NULL`)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to fetch pending migrations")
		return
	}

	// Step 2: Iterate over each migration
	for _, migration := range migrations {
		var keys []string
		err := d.db.Select(&keys, `SELECT key FROM meta_migration WHERE migration_id = ?`, migration.ID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to fetch keys for migration")
			continue
		}

		// Step 2.1: Check if there are at least 100 records
		if len(keys) < minKeysToMigrate {
			continue
		}

		// Step 2.2: Update migration_started_at to current timestamp
		_, err = d.db.Exec(`UPDATE migration SET migration_started_at = ? WHERE id = ?`, time.Now(), migration.ID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to update migration start time")
			continue
		}

		// Step 2.3: Retrieve data for keys
		values, _, err := retrieveBatchValues(ctx, d.p2pDataStore, keys, false, Store{})
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to retrieve batch values")
			continue
		}

		if err := d.uploadInBatches(ctx, keys, values); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to migrate some data")
			continue
		}

	}
}

func (d *MigrationMetaStore) uploadInBatches(ctx context.Context, keys []string, values [][]byte) error {
	const batchSize = 1000
	totalItems := len(keys)
	batches := (totalItems + batchSize - 1) / batchSize // Calculate number of batches needed
	var lastError error

	for i := 0; i < batches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchKeys := keys[start:end]
		batchData := values[start:end]

		uploadedKeys, err := d.cloud.UploadBatch(batchKeys, batchData)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Failed to upload batch %d", i+1)
			lastError = err
			continue
		}

		if err := batchDeleteRecords(d.p2pDataStore, uploadedKeys); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Failed to delete batch %d", i+1)
			lastError = err
			continue
		}

		log.WithContext(ctx).Infof("Successfully uploaded and deleted records for batch %d of %d", i+1, batches)
	}

	return lastError
}

type MetaStoreInterface interface {
	GetCountOfStaleData(ctx context.Context, staleTime time.Time) (int, error)
	GetStaleDataInBatches(ctx context.Context, batchSize, batchNumber int, duration time.Time) ([]string, error)
	GetPendingMigrationID(ctx context.Context) (int, error)
	CreateNewMigration(ctx context.Context) (int, error)
	InsertMetaMigrationData(ctx context.Context, migrationID int, keys []string) error
}

// GetCountOfStaleData returns the count of stale data where last_accessed is 3 months before.
func (d *MigrationMetaStore) GetCountOfStaleData(ctx context.Context, staleTime time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM meta WHERE last_accessed < ?`

	err := d.db.GetContext(ctx, &count, query, staleTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get count of stale data: %w", err)
	}
	return count, nil
}

// GetStaleDataInBatches retrieves stale data entries in batches from the meta table.
func (d *MigrationMetaStore) GetStaleDataInBatches(ctx context.Context, batchSize, batchNumber int, duration time.Time) ([]string, error) {
	offset := batchNumber * batchSize

	query := `
        SELECT key 
        FROM meta 
        WHERE last_accessed < ? 
        LIMIT ? OFFSET ?
    `

	rows, err := d.db.QueryxContext(ctx, query, duration, batchSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get stale data in batches: %w", err)
	}
	defer rows.Close()

	var staleData []string
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		staleData = append(staleData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}
	return staleData, nil
}

func (d *MigrationMetaStore) GetPendingMigrationID(ctx context.Context) (int, error) {
	var migrationID int
	query := `SELECT id FROM migration WHERE migration_started_at IS NULL LIMIT 1`

	err := d.db.GetContext(ctx, &migrationID, query)
	if err == sql.ErrNoRows {
		return 0, nil // No pending migrations
	} else if err != nil {
		return 0, fmt.Errorf("failed to get pending migration ID: %w", err)
	}

	return migrationID, nil
}

func (d *MigrationMetaStore) CreateNewMigration(ctx context.Context) (int, error) {
	query := `INSERT INTO migration (created_at, updated_at) VALUES (?, ?)`
	now := time.Now()

	result, err := d.db.ExecContext(ctx, query, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create new migration: %w", err)
	}

	migrationID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(migrationID), nil
}

func (d *MigrationMetaStore) InsertMetaMigrationData(ctx context.Context, migrationID int, keys []string) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt, err := tx.Preparex(`INSERT INTO meta_migration (key, migration_id, created_at, updated_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, key := range keys {
		if _, err := stmt.Exec(key, migrationID, now, now); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert meta migration data: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
