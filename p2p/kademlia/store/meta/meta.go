package meta

import (
	"context"
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
)

// Exponential backoff parameters
var (
	checkpointInterval = 5 * time.Second // Checkpoint interval in seconds
	dbName             = "data001-meta.sqlite3"
)

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

// DisabledKey represents the disabled key
type DisabledKey struct {
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"createdAt"`
}

// toDomain converts DisabledKey to domain.DisabledKey
func (r *DisabledKey) toDomain() (domain.DisabledKey, error) {
	return domain.DisabledKey{
		Key:       r.Key,
		CreatedAt: r.CreatedAt,
	}, nil
}

// Job represents the job to be run
type Job struct {
	JobType string // Insert, Delete
	Key     []byte

	TaskID string
	ReqID  string
}

// NewStore returns a new store
func NewStore(ctx context.Context, dataDir string) (*Store, error) {
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
	db.SetMaxOpenConns(250) // set appropriate value
	db.SetMaxIdleConns(10)  // set appropriate value

	s := &Store{
		db:     db,
		worker: worker,
	}

	if !s.checkStore() {
		if err = s.migrate(); err != nil {
			return nil, fmt.Errorf("cannot create table(s) in sqlite database: %w", err)
		}
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-120000;",
		"PRAGMA busy_timeout=120000;",
		"PRAGMA journal_size_limit=5242880;",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("cannot set sqlite database parameter: %w", err)
		}
	}

	s.db = db

	go s.start(ctx)
	// Run WAL checkpoint worker every 5 seconds
	go s.startCheckpointWorker(ctx)

	return s, nil
}

func (s *Store) checkStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='disabled_keys'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS disabled_keys(
        key TEXT PRIMARY KEY,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'disabled_keys': %w", err)
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

// Start method starts the run loop for the worker
func (s *Store) start(ctx context.Context) {
	for {
		select {
		case job := <-s.worker.JobQueue:
			if err := s.performJob(job); err != nil {
				log.WithError(err).Error("Failed to perform job")
			}
		case <-s.worker.quit:
			log.Info("exit sqlite meta db worker - quit signal received")
			return
		case <-ctx.Done():
			log.Info("exit sqlite meta db worker- ctx done signal received")
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
func (s *Store) Store(ctx context.Context, key []byte) error {
	job := Job{
		JobType: "Insert",
		Key:     key,
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

// PerformJob performs the job in the JobQueue
func (s *Store) performJob(j Job) error {
	switch j.JobType {
	case "Insert":
		err := s.storeDisabledKey(j.Key)
		if err != nil {
			log.WithError(err).WithField("taskID", j.TaskID).WithField("id", j.ReqID).Error("failed to store disable key record")
			return fmt.Errorf("failed to store disable key record: %w", err)
		}
	case "Delete":
		s.deleteRecord(j.Key)
	}

	return nil
}

// Checkpoint method for the store
func (s *Store) checkpoint() error {
	_, err := s.db.Exec("PRAGMA wal_checkpoint;")
	if err != nil {
		return fmt.Errorf("failed to checkpoint: %w", err)
	}
	return nil
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

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(_ context.Context, key string) error {
	r := DisabledKey{}
	err := s.db.Get(&r, `SELECT * FROM disabled_keys WHERE key = ?`, key)
	if err != nil {
		return fmt.Errorf("failed to get disabled record by key %s: %w", key, err)
	}

	if len(r.Key) == 0 {
		return fmt.Errorf("failed to get disabled record by key %s", key)
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

// deleteRecord a key/value pair from the Store
func (s *Store) deleteRecord(key []byte) {
	hkey := hex.EncodeToString(key)

	res, err := s.db.Exec("DELETE FROM disabled_keys WHERE key = ?", hkey)
	if err != nil {
		log.P2P().Debugf("cannot delete disabled keys record by key %s: %v", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		log.P2P().Debugf("failed to delete disabled key record by key %s: %v", hkey, err)
	} else if rowsAffected == 0 {
		log.P2P().Debugf("failed to delete disabled key record by key %s", hkey)
	}
}

func (s *Store) storeDisabledKey(key []byte) error {
	operation := func() error {
		hkey := hex.EncodeToString(key)

		now := time.Now().UTC()
		r := DisabledKey{Key: hkey, CreatedAt: now}
		res, err := s.db.NamedExec(`INSERT INTO disabled_keys(key, createdAt) values(:key, :createdat) ON CONFLICT(key) DO UPDATE SET createdAt=:createdAt`, r)
		if err != nil {
			return fmt.Errorf("cannot insert or update disabled key record with key %s: %w", hkey, err)
		}

		if rowsAffected, err := res.RowsAffected(); err != nil {
			return fmt.Errorf("failed to insert/update disabled key record with key %s: %w", hkey, err)
		} else if rowsAffected == 0 {
			return fmt.Errorf("failed to insert/update disabled key record with key %s", hkey)
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

// GetDisabledKeys returns all disabled keys
func (s *Store) GetDisabledKeys(from time.Time) (retKeys domain.DisabledKeys, err error) {
	var keys []DisabledKey
	if err := s.db.Select(&keys, "SELECT * FROM disabled_keys where createdAt > from", from); err != nil {
		return nil, fmt.Errorf("error reading disabled keys from database: %w", err)
	}

	retKeys = make(domain.DisabledKeys, 0, len(keys))

	for _, key := range keys {
		domainKey, err := key.toDomain()
		if err != nil {
			return nil, fmt.Errorf("error converting disabled key to domain: %w", err)
		}
		retKeys = append(retKeys, domainKey)
	}

	return retKeys, nil
}
