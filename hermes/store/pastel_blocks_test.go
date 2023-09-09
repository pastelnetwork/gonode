package store

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text PRIMARY KEY, path_to_art_image_file text, datetime_fingerprint_added_to_database text, thumbnail_of_image text, request_type text, open_api_group_id_string text, collection_name_string text, registration_ticket_txid text, txid_timestamp integer)`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestSQLiteStore_FetchAllTxIDs(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("setupDatabase() err = %v; want %v", err, nil)
	}

	store := &SQLiteStore{db: db}

	tests := []struct {
		name        string
		insertQuery string
		want        map[string]bool
	}{
		{
			name:        "no txid",
			insertQuery: "",
			want:        map[string]bool{},
		},
		{
			name:        "single txid",
			insertQuery: "INSERT INTO image_hash_to_image_fingerprint_table (registration_ticket_txid) VALUES ('txid1234')",
			want:        map[string]bool{"txid1234": true},
		},
		{
			name:        "multiple txids",
			insertQuery: "INSERT INTO image_hash_to_image_fingerprint_table (registration_ticket_txid) VALUES ('txid1234'), ('txid5678')",
			want:        map[string]bool{"txid1234": true, "txid5678": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.insertQuery != "" {
				_, err := db.Exec(tt.insertQuery)
				if err != nil {
					t.Fatalf("Failed to insert data: %v", err)
				}
			}

			got, err := store.FetchAllTxIDs()
			if err != nil {
				t.Errorf("FetchAllTxIDs() err = %v; want %v", err, nil)
			}
			if len(got) != len(tt.want) {
				t.Errorf("FetchAllTxIDs() got = %v; want %v", got, tt.want)
			}
			for k := range tt.want {
				if _, ok := got[k]; !ok {
					t.Errorf("Expected key %v was not found in got", k)
				}
			}
		})
	}
}

func TestSQLiteStore_UpdateTxIDTimestamp(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("setupDatabase() err = %v; want %v", err, nil)
	}

	store := &SQLiteStore{db: db}

	tests := []struct {
		name             string
		initialDataQuery string
		updateTxID       string
		wantError        bool
	}{
		{
			name:             "txid exists",
			initialDataQuery: "INSERT INTO image_hash_to_image_fingerprint_table (registration_ticket_txid, txid_timestamp) VALUES ('txid1234', 100)",
			updateTxID:       "txid1234",
			wantError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.initialDataQuery != "" {
				_, err := db.Exec(tt.initialDataQuery)
				if err != nil {
					t.Fatalf("Failed to insert data: %v", err)
				}
			}

			err := store.UpdateTxIDTimestamp(tt.updateTxID)
			if (err != nil) != tt.wantError {
				t.Errorf("UpdateTxIDTimestamp() err = %v; wantError = %v", err, tt.wantError)
				return
			}

			// Check that the timestamp is set to 0
			var timestamp int
			row := db.QueryRow("SELECT txid_timestamp FROM image_hash_to_image_fingerprint_table WHERE registration_ticket_txid = ?", tt.updateTxID)
			err = row.Scan(&timestamp)
			if err != nil {
				t.Errorf("Failed to fetch updated timestamp: %v", err)
				return
			}

			if timestamp != 0 {
				t.Errorf("Timestamp was not updated to 0. Got: %d; want: 0", timestamp)
			}
		})
	}
}
