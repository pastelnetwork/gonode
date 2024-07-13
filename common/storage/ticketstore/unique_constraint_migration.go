package ticketstore

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func removePrimaryKeyMigration(db sqlx.DB) error {
	var columns []struct {
		Cid        int     `db:"cid"`
		Name       string  `db:"name"`
		Type       string  `db:"type"`
		Notnull    int     `db:"notnull"`
		Dflt_value *string `db:"dflt_value"`
		Pk         int     `db:"pk"`
	}

	err := db.Select(&columns, "PRAGMA table_info(files);")
	if err != nil {
		return err
	}

	primaryKeyColumns := []string{}
	for _, column := range columns {
		if column.Pk == 1 || column.Pk == 2 {
			primaryKeyColumns = append(primaryKeyColumns, column.Name)
		}
	}

	// Check if the primary key is only on file_id
	if len(primaryKeyColumns) == 1 && primaryKeyColumns[0] == "file_id" {
		tx := db.MustBegin()

		// Create the new table with the desired schema
		tx.MustExec(`
		CREATE TABLE IF NOT EXISTS files_new (
			file_id TEXT NOT NULL,
			upload_timestamp DATETIME,
			path TEXT,
			file_index TEXT,
			base_file_id TEXT,
			task_id TEXT,
			reg_txid TEXT,
			activation_txid TEXT,
			req_burn_txn_amount FLOAT,
			burn_txn_id TEXT,
			req_amount FLOAT,
			is_concluded BOOLEAN,
			cascade_metadata_ticket_id TEXT,
			uuid_key TEXT,
			hash_of_original_big_file TEXT,
			name_of_original_big_file_with_ext TEXT,
			size_of_original_big_file FLOAT,
			data_type_of_original_big_file TEXT,
			start_block INTEGER,
			done_block INTEGER,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			pastel_id TEXT,
			passphrase TEXT,
			PRIMARY KEY (file_id, base_file_id)
		);
		`)

		// Copy data from the old table to the new table
		tx.MustExec(`INSERT INTO files_new SELECT * FROM files;`)

		// Drop the old table
		tx.MustExec(`DROP TABLE files;`)

		// Rename the new table to the old table name
		tx.MustExec(`ALTER TABLE files_new RENAME TO files;`)

		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("primary-key update migration executed successfully")
	} else {
		fmt.Println("primary key constraint already existed on file_id & base_file_id. Not executing the migration")
	}

	return nil
}
