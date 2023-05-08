package sqlite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"

	"github.com/pastelnetwork/gonode/common/log"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/utils"
)

// Exponential backoff parameters
var (
	checkpointInterval = 5 * time.Second // Checkpoint interval in seconds
	//dbLock             sync.Mutex
	replicateInterval = 10 * time.Minute
	dbName            = "data001.sqlite3"
	dbFilePath        = ""
)

// Job represents the job to be run
type Job struct {
	JobType      string // Insert, Update or Delete
	Key          []byte
	Value        []byte
	Values       [][]byte
	ReplicatedAt time.Time
	TaskID       string
	ReqID        string
}

// Worker represents the worker that executes the job
type Worker struct {
	JobQueue chan Job
	quit     chan bool
}

// Store is the main struct
type Store struct {
	db     *sqlx.DB
	worker *Worker
}

// Record is a data record
type Record struct {
	Key          string
	Data         []byte
	IsOriginal   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ReplicatedAt time.Time
}

// NewStore returns a new store
func NewStore(ctx context.Context, dataDir string, replicate time.Duration, _ time.Duration) (*Store, error) {
	worker := &Worker{
		JobQueue: make(chan Job, 500),
		quit:     make(chan bool),
	}

	log.P2P().WithContext(ctx).Infof("p2p data dir: %v", dataDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return nil, fmt.Errorf("mkdir %q: %w", dataDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot create data folder: %w", err)
	}

	dbFile := path.Join(dataDir, dbName)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}

	s := &Store{
		worker: worker,
		db:     db,
	}

	if !s.checkStore() {
		if err = s.migrate(); err != nil {
			return nil, fmt.Errorf("cannot create table(s) in sqlite database: %w", err)
		}
	}

	if !s.checkReplicateStore() {
		if err = s.migrateReplication(); err != nil {
			return nil, fmt.Errorf("cannot create table(s) in sqlite database: %w", err)
		}
	}

	_, err = db.Exec(`PRAGMA journal_mode=WAL;
	PRAGMA synchronous=NORMAL;
	PRAGMA cache_size=-20000;
	PRAGMA journal_size_limit=5242880`)
	if err != nil {
		return nil, fmt.Errorf("cannot set sqlite database parameters: %w", err)
	}

	s.db = db
	replicateInterval = replicate
	dbFilePath = dbFile

	go s.start(ctx)
	// Run WAL checkpoint worker every 5 seconds
	go s.startCheckpointWorker(ctx)

	return s, nil
}

func (s *Store) checkStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='data'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) checkReplicateStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='replication_info'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS data(
        key TEXT PRIMARY KEY,
        data BLOB NOT NULL,
        is_original BOOL DEFAULT FALSE,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        replicatedAt DATETIME,
        republishedAt DATETIME
    );
    `

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'data': %w", err)
	}

	return nil
}

func (s *Store) migrateReplication() error {
	replicateQuery := `
    CREATE TABLE IF NOT EXISTS replication_info(
        id TEXT PRIMARY KEY,
        ip TEXT NOT NULL,
		port INTEGER NOT NULL,
        is_active BOOL DEFAULT FALSE,
		is_adjusted BOOL DEFAULT FALSE,
        lastReplicatedAt DATETIME,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `

	if _, err := s.db.Exec(replicateQuery); err != nil {
		return fmt.Errorf("failed to create table 'replication_info': %w", err)
	}

	return nil
}

func (s *Store) startCheckpointWorker(ctx context.Context) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 1 * time.Minute
	b.InitialInterval = 100 * time.Millisecond

	for {
		err := backoff.RetryNotify(func() error {
			err := s.checkpoint()
			if err == nil {
				// If no error, delay for 5 seconds.
				time.Sleep(checkpointInterval)
			}
			return err
		}, b, func(err error, duration time.Duration) {
			log.WithContext(ctx).WithField("duration", duration).Error("Failed to perform checkpoint, retrying...")
		})

		if err == nil {
			b.Reset()
			b.MaxElapsedTime = 1 * time.Minute
			b.InitialInterval = 100 * time.Millisecond
		}

		select {
		case <-ctx.Done():
			log.WithContext(ctx).Info("Stopping checkpoint worker because of context cancel")
			return
		case <-s.worker.quit:
			log.WithContext(ctx).Info("Stopping checkpoint worker because of quit signal")
			return
		default:
		}
	}
}

type nodeReplicationInfo struct {
	LastReplicated *time.Time `db:"lastReplicatedAt"`
	UpdatedAt      time.Time  `db:"updatedAt"`
	CreatedAt      time.Time  `db:"createdAt"`
	Active         bool       `db:"is_active"`
	Adjusted       bool       `db:"is_adjusted"`
	IP             string     `db:"ip"`
	Port           int        `db:"port"`
	ID             string     `db:"id"`
}

func (n *nodeReplicationInfo) toDomain() domain.NodeReplicationInfo {
	return domain.NodeReplicationInfo{
		LastReplicatedAt: n.LastReplicated,
		UpdatedAt:        n.UpdatedAt,
		CreatedAt:        n.CreatedAt,
		Active:           n.Active,
		IsAdjusted:       n.Adjusted,
		IP:               n.IP,
		Port:             n.Port,
		ID:               []byte(n.ID),
	}
}

// Start method starts the run loop for the worker
func (s *Store) start(ctx context.Context) {
	for {
		select {
		case job := <-s.worker.JobQueue:
			if err := s.performJob(job); err != nil {
				log.WithError(err).Error("Failed to perform job")
			}
		case <-s.worker.quit:
			log.Info("exit sqlite db worker - quit signal received")
			return
		case <-ctx.Done():
			log.Info("exit sqlite db worker- ctx done signal received")
			return
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// Store function creates a new job and pushes it into the JobQueue
func (s *Store) Store(ctx context.Context, key []byte, value []byte) error {

	job := Job{
		JobType: "Insert",
		Key:     key,
		Value:   value,
	}

	if val := ctx.Value(log.TaskIDKey); val != nil {
		switch val := val.(type) {
		case string:
			job.TaskID = val
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.worker.JobQueue <- job:
	}

	return nil
}

// StoreBatch stores a batch of key/value pairs for the local node with the replication
func (s *Store) StoreBatch(ctx context.Context, values [][]byte) error {
	job := Job{
		JobType: "BatchInsert",
		Values:  values,
	}

	if val := ctx.Value(log.TaskIDKey); val != nil {
		switch val := val.(type) {
		case string:
			job.TaskID = val
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.worker.JobQueue <- job:
	}

	return nil
}

// Delete a key/value pair from the store
func (s *Store) Delete(_ context.Context, key []byte) {
	job := Job{
		JobType: "Delete",
		Key:     key,
	}

	s.worker.JobQueue <- job
}

// DeleteAll the records in store
func (s *Store) DeleteAll(ctx context.Context) error {
	job := Job{
		JobType: "DeleteAll",
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.worker.JobQueue <- job:
	}

	return nil
}

// UpdateKeyReplication updates the replication status of the key
func (s *Store) UpdateKeyReplication(ctx context.Context, key []byte) error {
	job := Job{
		JobType: "Update",
		Key:     key,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.worker.JobQueue <- job:
	}

	return nil
}

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(_ context.Context, key []byte) ([]byte, error) {
	hkey := hex.EncodeToString(key)

	r := Record{}
	err := s.db.Get(&r, `SELECT data FROM data WHERE key = ?`, hkey)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by key %s: %w", hkey, err)
	}

	return r.Data, nil
}

// Checkpoint method for the store
func (s *Store) checkpoint() error {
	//dbLock.Lock()
	//defer dbLock.Unlock()

	_, err := s.db.Exec("PRAGMA wal_checkpoint;")
	if err != nil {
		return fmt.Errorf("failed to checkpoint: %w", err)
	}
	return nil
}

// PerformJob performs the job in the JobQueue
func (s *Store) performJob(j Job) error {
	switch j.JobType {
	case "Insert":
		err := s.storeRecord(j.Key, j.Value)
		if err != nil {
			log.WithError(err).WithField("taskID", j.TaskID).WithField("id", j.ReqID).Error("failed to store record")
			return fmt.Errorf("failed to store record: %w", err)
		}

	case "BatchInsert":
		err := s.storeBatchRecord(j.Values)
		if err != nil {
			log.WithError(err).WithField("taskID", j.TaskID).WithField("id", j.ReqID).Error("failed to store batch records")
			return fmt.Errorf("failed to store batch record: %w", err)
		}

		log.WithField("taskID", j.TaskID).WithField("id", j.ReqID).Info("successfully stored batch records")
	case "Update":
		err := s.updateKeyReplication(j.Key, j.ReplicatedAt)
		if err != nil {
			return fmt.Errorf("failed to update key replication: %w", err)
		}
	case "Delete":
		s.deleteRecord(j.Key)
	case "DeleteAll":
		err := s.deleteAll()
		if err != nil {
			return fmt.Errorf("failed to delete record: %w", err)
		}
	}

	return nil
}

// storeRecord will store a key/value pair for the local node
func (s *Store) storeRecord(key []byte, value []byte) error {

	operation := func() error {
		hkey := hex.EncodeToString(key)

		now := time.Now().UTC()
		r := Record{Key: hkey, Data: value, UpdatedAt: now}
		res, err := s.db.NamedExec(`INSERT INTO data(key, data) values(:key, :data) ON CONFLICT(key) DO UPDATE SET data=:data,updatedAt=:updatedat`, r)
		if err != nil {
			return fmt.Errorf("cannot insert or update record with key %s: %w", hkey, err)
		}

		if rowsAffected, err := res.RowsAffected(); err != nil {
			return fmt.Errorf("failed to insert/update record with key %s: %w", hkey, err)
		} else if rowsAffected == 0 {
			return fmt.Errorf("failed to insert/update record with key %s", hkey)
		}

		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(operation, b)
	if err != nil {
		return fmt.Errorf("error storing data: %w", err)
	}

	return nil
}

// storeBatchRecord will store a batch of values with their SHA256 hash as the key
func (s *Store) storeBatchRecord(values [][]byte) error {
	operation := func() error {
		tx, err := s.db.Beginx()
		if err != nil {
			return fmt.Errorf("cannot begin transaction: %w", err)
		}

		// Prepare insert statement
		stmt, err := tx.PrepareNamed(`INSERT INTO data(key, data) values(:key, :data) ON CONFLICT(key) DO UPDATE SET data=:data,updatedAt=:updatedat`)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("cannot prepare statement: %w", err)
		}

		// For each value, calculate its hash and insert into DB
		now := time.Now().UTC()
		for i := 0; i < len(values); i++ {
			// Compute the SHA256 hash
			h := sha256.New()
			h.Write(values[i])
			hashed, err := utils.Sha3256hash(values[i])
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("cannot compute hash: %w", err)
			}

			hkey := hex.EncodeToString(hashed)
			r := Record{Key: hkey, Data: values[i], UpdatedAt: now}

			// Execute the insert statement
			_, err = stmt.Exec(r)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("cannot insert or update record with key %s: %w", hkey, err)
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("cannot commit transaction: %w", err)
		}

		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(operation, b)
	if err != nil {
		return fmt.Errorf("error storing data: %w", err)
	}

	return nil
}

// deleteRecord a key/value pair from the Store
func (s *Store) deleteRecord(key []byte) {
	hkey := hex.EncodeToString(key)

	res, err := s.db.Exec("DELETE FROM data WHERE key = ?", hkey)
	if err != nil {
		log.P2P().Debugf("cannot delete record by key %s: %v", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		log.P2P().Debugf("failed to delete record by key %s: %v", hkey, err)
	} else if rowsAffected == 0 {
		log.P2P().Debugf("failed to delete record by key %s", hkey)
	}
}

// updateKeyReplication updates the replication time for a key
func (s *Store) updateKeyReplication(key []byte, replicatedAt time.Time) error {
	keyStr := hex.EncodeToString(key)
	_, err := s.db.Exec(`UPDATE data SET replicatedAt = ? WHERE key = ?`, replicatedAt, keyStr)
	if err != nil {
		return fmt.Errorf("failed to update replicated records: %v", err)
	}

	return err
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(ctx context.Context, from time.Time) (keys [][]byte) {
	if err := s.db.Select(&keys, `SELECT key FROM data WHERE createdAt > ?`, from); err != nil {
		log.P2P().WithError(err).WithContext(ctx).Errorf("failed to get records for replication older than %s", from)
		return nil
	}

	return keys
}

func (s *Store) deleteAll() error {
	res, err := s.db.Exec("DELETE FROM data")
	if err != nil {
		return fmt.Errorf("cannot delete ALL records: %w", err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("failed to delete ALL records: %w", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("failed to delete ALL records")
	}

	return nil
}

// Count the records in store
func (s *Store) Count(_ context.Context) (int, error) {
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM data`)
	if err != nil {
		return -1, fmt.Errorf("failed to get count of records: %w", err)
	}

	return count, nil
}

// Stats returns stats of store
func (s *Store) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	fi, err := os.Stat(dbFilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get p2p db size")
	} else {
		stats["p2p_db_size"] = utils.BytesToMB(uint64(fi.Size()))
	}

	if count, err := s.Count(ctx); err == nil {
		stats["p2p_db_records_count"] = count
	} else {
		log.WithContext(ctx).WithError(err).Error("failed to get p2p records count")
	}

	return stats, nil
}

// Close the store
func (s *Store) Close(ctx context.Context) {
	s.worker.Stop()

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.P2P().WithContext(ctx).Errorf("Failed to close database: %s", err)
		}
	}
}

// GetAllReplicationInfo returns all records in replication table
func (s *Store) GetAllReplicationInfo(_ context.Context) ([]domain.NodeReplicationInfo, error) {
	r := []nodeReplicationInfo{}
	err := s.db.Select(&r, "SELECT * FROM replication_info")
	if err != nil {
		return nil, fmt.Errorf("failed to get replication info %w", err)
	}

	list := []domain.NodeReplicationInfo{}
	for _, v := range r {
		list = append(list, v.toDomain())
	}

	return list, nil
}

// UpdateReplicationInfo updates replication info
func (s *Store) UpdateReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	_, err := s.db.Exec(`UPDATE replication_info SET ip = ?, is_active = ?, is_adjusted = ?, lastReplicatedAt = ?, updatedAt =?, port = ? WHERE id = ?`,
		rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port, string(rep.ID))
	if err != nil {
		return fmt.Errorf("failed to update replicated records: %v", err)
	}

	return err
}

// AddReplicationInfo adds replication info
func (s *Store) AddReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	_, err := s.db.Exec(`INSERT INTO replication_info(id, ip, is_active, is_adjusted, lastReplicatedAt, updatedAt, port) values(?, ?, ?, ?, ?, ?, ?)`,
		string(rep.ID), rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port)
	if err != nil {
		return fmt.Errorf("failed to update replicate record: %v", err)
	}

	return err
}
