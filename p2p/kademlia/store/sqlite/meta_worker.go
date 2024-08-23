package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	commitLastAccessedInterval = 90 * time.Second
	migrationExecutionTicker   = 12 * time.Hour
	migrationMetaDB            = "data001-migration-meta.sqlite3"
	accessUpdateBufferSize     = 100000
	commitInsertsInterval      = 10 * time.Second
	metaSyncBatchSize          = 5000
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

// UpdateMessage holds the key and the last accessed time.
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

	if err := setPragmas(db); err != nil {
		log.WithContext(ctx).WithError(err).Error("error executing pragmas")
	}

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

	go func() {
		if handler.isMetaSyncRequired() {
			err := handler.syncMetaWithData(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error syncing meta with p2p data")
			}
		}
	}()

	go handler.startLastAccessedUpdateWorker(ctx)
	go handler.startInsertWorker(ctx)
	go handler.startMigrationExecutionWorker(ctx)
	log.WithContext(ctx).Info("MigrationMetaStore workers started")

	return handler, nil
}

func setPragmas(db *sqlx.DB) error {
	// Set journal mode to WAL
	_, err := db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return err
	}
	// Set synchronous to NORMAL
	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		return err
	}

	// Set cache size
	_, err = db.Exec("PRAGMA cache_size=-262144;")
	if err != nil {
		return err
	}

	// Set busy timeout
	_, err = db.Exec("PRAGMA busy_timeout=5000;")
	if err != nil {
		return err
	}

	return nil
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
	var offset int

	query := `SELECT key, data, updatedAt FROM data LIMIT ? OFFSET ?`
	insertQuery := `
	INSERT INTO meta (key, last_accessed, access_count, data_size)
	VALUES (?, ?, 1, ?)
	ON CONFLICT(key) DO
	UPDATE SET
	last_accessed = EXCLUDED.last_accessed,
	    data_size = EXCLUDED.data_size,
		access_count = access_count + 1`

	for {
		rows, err := d.p2pDataStore.Queryx(query, metaSyncBatchSize, offset)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error querying p2p data store")
			return err
		}

		tx, err := d.db.Beginx()
		if err != nil {
			rows.Close()
			log.WithContext(ctx).WithError(err).Error("Failed to start transaction")
			return err
		}

		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			tx.Rollback()
			rows.Close()
			log.WithContext(ctx).WithError(err).Error("Failed to prepare statement")
			return err
		}

		var batchProcessed bool
		for rows.Next() {
			var r Record
			var t *time.Time

			if err := rows.Scan(&r.Key, &r.Data, &t); err != nil {
				log.WithContext(ctx).WithError(err).Error("Error scanning row from p2p data store")
				continue
			}
			if t != nil {
				r.UpdatedAt = *t
			}

			if _, err := stmt.Exec(r.Key, r.UpdatedAt, len(r.Data)); err != nil {
				log.WithContext(ctx).WithField("key", r.Key).WithError(err).Error("error inserting key to meta")
				continue
			}

			batchProcessed = true
		}

		stmt.Close()
		if err := rows.Err(); err != nil {
			tx.Rollback()
			rows.Close()
			log.WithContext(ctx).WithError(err).Error("Error iterating rows")
			return err
		}

		if batchProcessed {
			if err := tx.Commit(); err != nil {
				rows.Close()
				log.WithContext(ctx).WithError(err).Error("Failed to commit transaction")
				return err
			}
		} else {
			tx.Rollback()
			rows.Close()
			break
		}

		rows.Close()
		if !batchProcessed {
			tx.Rollback()
			log.WithContext(ctx).Info("no rows processed, rolling back and breaking.")
			break
		}
		offset += metaSyncBatchSize //
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
			LastAccessTime: time.Now().UTC(),
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
			return
		}
	}
}

func (d *MigrationMetaStore) commitLastAccessedUpdates(ctx context.Context) {
	keysToUpdate := make(map[string]time.Time)
	d.updates.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			log.WithContext(ctx).Error("Error converting key to string (commitLastAccessedUpdates)")
			return false
		}
		v, ok := value.(time.Time)
		if !ok {
			log.WithContext(ctx).Error("Error converting value to time.Time (commitLastAccessedUpdates)")
			return false
		}
		keysToUpdate[k] = v

		return true // continue
	})

	if len(keysToUpdate) == 0 {
		return
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error starting transaction (commitLastAccessedUpdates)")
		return
	}

	stmt, err := tx.Prepare(`
	INSERT INTO meta (key, last_accessed, access_count) 
	VALUES (?, ?, 1) 
	ON CONFLICT(key) DO 
	UPDATE SET 
		last_accessed = EXCLUDED.last_accessed,
		access_count = access_count + 1`)
	if err != nil {
		tx.Rollback() // Roll back the transaction on error
		log.WithContext(ctx).WithError(err).Error("Error preparing statement (commitLastAccessedUpdates)")
		return
	}
	defer stmt.Close()

	for k, v := range keysToUpdate {
		_, err := stmt.Exec(k, v)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("key", k).Error("Error executing statement (commitLastAccessedUpdates)")
		}
	}

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
	keysToUpdate := make(map[string]UpdateMessage)
	d.inserts.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			log.WithContext(ctx).Error("Error converting key to string (commitInserts)")
			return false
		}
		v, ok := value.(UpdateMessage)
		if !ok {
			log.WithContext(ctx).Error("Error converting value to UpdateMessage (commitInserts)")
			return false
		}

		keysToUpdate[k] = v

		return true // continue
	})

	if len(keysToUpdate) == 0 {
		return
	}

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

	for k, v := range keysToUpdate {
		accessCount := 1
		_, err := stmt.Exec(k, v.LastAccessTime, accessCount, v.Size)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("key", k).Error("Error executing statement (commitInserts)")
		}
	}

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

type Migration struct {
	ID                  int          `db:"id"`
	MigrationStartedAt  sql.NullTime `db:"migration_started_at"`
	MigrationFinishedAt sql.NullTime `db:"migration_finished_at"`
}
type Migrations []Migration

type MigrationKey struct {
	Key         string `db:"key"`
	MigrationID int    `db:"migration_id"`
	IsMigrated  bool   `db:"is_migrated"`
}

type MigrationKeys []MigrationKey

func (d *MigrationMetaStore) checkAndExecuteMigration(ctx context.Context) {
	// Check the available disk space
	isLow, err := utils.CheckDiskSpace(lowSpaceThresholdGB)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("migration worker: check disk space failed")
	}

	//if !isLow {
	// Disk space is sufficient, stop migration
	//return
	//}

	log.WithContext(ctx).WithField("islow", isLow).Info("Starting data migration")
	// Step 1: Fetch pending migrations
	var migrations Migrations

	err = d.db.Select(&migrations, `SELECT id FROM migration WHERE migration_started_at IS NULL or migration_finished_at IS NULL`)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to fetch pending migrations")
		return
	}
	log.WithContext(ctx).WithField("count", len(migrations)).Info("Fetched pending migrations")

	// Iterate over each migration
	for _, migration := range migrations {
		log.WithContext(ctx).WithField("migration_id", migration.ID).Info("Processing migration")

		if err := d.ProcessMigrationInBatches(ctx, migration); err != nil {
			log.WithContext(ctx).WithError(err).WithField("migration_id", migration.ID).Error("Failed to process migration")
			continue
		}
	}
}

func (d *MigrationMetaStore) ProcessMigrationInBatches(ctx context.Context, migration Migration) error {
	var migKeys MigrationKeys
	err := d.db.Select(&migKeys, `SELECT key FROM meta_migration WHERE migration_id = ?`, migration.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch keys for migration %d: %w", migration.ID, err)
	}

	totalKeys := len(migKeys)
	if totalKeys == 0 {
		return nil
	}

	//if totalKeys < minKeysToMigrate {
	//	log.WithContext(ctx).WithField("migration_id", migration.ID).WithField("keys-count", totalKeys).Info("Skipping migration due to insufficient keys")
	//	return nil
	//}

	migratedKeys := 0
	var keys []string
	for _, key := range migKeys {
		if key.IsMigrated {
			migratedKeys++
		} else {
			keys = append(keys, key.Key)
		}
	}

	nonMigratedKeys := len(keys)
	maxBatchSize := 5000
	for start := 0; start < nonMigratedKeys; start += maxBatchSize {
		end := start + maxBatchSize
		if end > nonMigratedKeys {
			end = nonMigratedKeys
		}

		batchKeys := keys[start:end]

		// Mark migration as started
		if start == 0 && !migration.MigrationStartedAt.Valid {
			if _, err := d.db.Exec(`UPDATE migration SET migration_started_at = ? WHERE id = ?`, time.Now(), migration.ID); err != nil {
				return err
			}
		}

		// Retrieve and upload data for current batch
		if err := d.processSingleBatch(ctx, batchKeys); err != nil {
			log.WithContext(ctx).WithError(err).WithField("migration_id", migration.ID).WithField("batch-keys-count", len(batchKeys)).Error("Failed to process batch")
			return fmt.Errorf("failed to process batch: %w - exiting now", err)
		}
	}

	// Mark migration as finished if all keys are migrated
	var leftMigKeys MigrationKeys
	err = d.db.Select(&leftMigKeys, `SELECT key FROM meta_migration WHERE migration_id = ? and is_migrated = ?`, migration.ID, false)
	if err != nil {
		return fmt.Errorf("failed to fetch keys for migration %d: %w", migration.ID, err)
	}

	if len(leftMigKeys) == 0 {
		if _, err := d.db.Exec(`UPDATE migration SET migration_finished_at = ? WHERE id = ?`, time.Now(), migration.ID); err != nil {
			return fmt.Errorf("failed to mark migration %d as finished: %w", migration.ID, err)
		}
	}

	log.WithContext(ctx).WithField("migration_id", migration.ID).WithField("tota-keys-count", totalKeys).WithField("migrated_in_current_iteration", nonMigratedKeys).Info("Migration processed successfully")

	return nil
}

func (d *MigrationMetaStore) processSingleBatch(ctx context.Context, keys []string) error {
	values, nkeys, err := d.retrieveBatchValuesToMigrate(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to retrieve batch values: %w", err)
	}

	if err := d.uploadInBatches(ctx, nkeys, values); err != nil {
		return fmt.Errorf("failed to upload batch values: %w", err)
	}

	return nil
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

		if err := batchSetMigratedRecords(d.p2pDataStore, uploadedKeys); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Failed to delete batch %d", i+1)
			lastError = err
			continue
		}

		if err := d.batchSetMigrated(uploadedKeys); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Failed to batch is_migrated %d", i+1)
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
	query := `SELECT COUNT(*)
FROM meta m
LEFT JOIN (
    SELECT DISTINCT mm.key
    FROM meta_migration mm
    JOIN migration mg ON mm.migration_id = mg.id
    WHERE mg.migration_started_at > $1
    AND mm.is_migrated = true   
       OR mm.created_at > $1
) AS recent_migrations ON m.key = recent_migrations.key
WHERE m.last_accessed < $1
AND recent_migrations.key IS NULL;`

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
SELECT m.key
FROM meta m
LEFT JOIN (
    SELECT DISTINCT mm.key
    FROM meta_migration mm
    JOIN migration mg ON mm.migration_id = mg.id
    WHERE mg.migration_started_at > $1 
    AND mm.is_migrated = true   
       OR mm.created_at > $1
) AS recent_migrations ON m.key = recent_migrations.key
WHERE m.last_accessed < $1
AND recent_migrations.key IS NULL
LIMIT $2 OFFSET $3;`

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

func (d *MigrationMetaStore) batchSetMigrated(keys []string) error {
	if len(keys) == 0 {
		log.P2P().Info("no keys provided for batch update (is_migrated)")
		return nil
	}

	// Create a parameter string for SQL query (?, ?, ?, ...)
	paramStr := strings.Repeat("?,", len(keys)-1) + "?"

	// Create the SQL statement
	query := fmt.Sprintf("UPDATE meta_migration set is_migrated = true WHERE key IN (%s)", paramStr)

	// Execute the query
	res, err := d.db.Exec(query, stringArgsToInterface(keys)...)
	if err != nil {
		return fmt.Errorf("cannot batch update records (is_migrated): %w", err)
	}

	// Optionally check rows affected
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("failed to get rows affected for batch update(is_migrated): %w", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no rows affected for batch update (is_migrated)")
	}

	return nil
}

func (d *MigrationMetaStore) retrieveBatchValuesToMigrate(ctx context.Context, keys []string) ([][]byte, []string, error) {
	if len(keys) == 0 {
		return [][]byte{}, []string{}, nil // Return empty if no keys provided
	}

	placeholders := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		placeholders[i] = "?"
		args[i] = key
	}

	query := fmt.Sprintf(`SELECT key, data, is_on_cloud FROM data WHERE key IN (%s) AND is_on_cloud = false`, strings.Join(placeholders, ","))
	rows, err := d.p2pDataStore.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve records: %w", err)
	}
	defer rows.Close()

	// Assume pre-allocation for efficiency
	values := make([][]byte, 0, len(keys))
	foundKeys := make([]string, 0, len(keys))

	for rows.Next() {
		var key string
		var value []byte
		var isOnCloud bool
		if err := rows.Scan(&key, &value, &isOnCloud); err != nil {
			return nil, nil, fmt.Errorf("failed to scan key and value: %w", err)
		}
		if len(value) > 1 && !isOnCloud {
			values = append(values, value)
			foundKeys = append(foundKeys, key)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, foundKeys, fmt.Errorf("rows processing error: %w", err)
	}

	return values, foundKeys, nil
}
