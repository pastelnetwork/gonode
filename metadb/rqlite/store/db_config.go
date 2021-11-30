package store

// DBConfig represents the configuration of the underlying SQLite database.
type DBConfig struct {
	// Whether the database is in-memory only.
	Memory bool `json:"memory"`

	// SQLite on-disk path
	OnDiskPath string `json:"on_disk_path,omitempty"`

	// Enforce Foreign Key constraints
	FKConstraints bool `json:"fk_constraints"`
}

// NewDBConfig returns a new DB config instance.
func NewDBConfig(memory bool) *DBConfig {
	return &DBConfig{Memory: memory}
}
