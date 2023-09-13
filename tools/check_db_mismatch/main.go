package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
)

func fetchTxIDsFromDB(db *sqlx.DB) (map[string]bool, error) {
	txIDs := make(map[string]bool)

	rows, err := db.Query("SELECT registration_ticket_txid FROM image_hash_to_image_fingerprint_table WHERE registration_ticket_txid IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txID string
	for rows.Next() {
		if err := rows.Scan(&txID); err != nil {
			return nil, err
		}
		txIDs[txID] = true
	}

	return txIDs, rows.Err()
}

func main() {

	db1, err := sqlx.Connect("sqlite3", "hermes-154.12.235.12.sqlite")
	if err != nil {
		log.Errorf("cannot open dd-service database: %w", err)
		return
	}

	defer db1.Close()

	db2, err := sqlx.Connect("sqlite3", "hermes-3.19.48.187.sqlite")
	if err != nil {
		log.Errorf("cannot open dd-service database: %w", err)
		return
	}

	defer db2.Close()

	txIDs1, err := fetchTxIDsFromDB(db1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of txIDs in hermes 1", len(txIDs1))

	txIDs2, err := fetchTxIDsFromDB(db2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of txIDs in hermes 2", len(txIDs2))

	fmt.Println("TxIDs in hermes-154.12.235.16.sqlite but not in hermes-154.38.161.58.sqlite:")
	for txID := range txIDs1 {
		if _, exists := txIDs2[txID]; !exists {
			fmt.Println(txID)
		}
	}

	fmt.Println("\nTxIDs in hermes-154.12.235.41.sqlite but not in hermes-154.12.235.41.sqlite:")
	for txID := range txIDs2 {
		if _, exists := txIDs1[txID]; !exists {
			fmt.Println(txID)
		}
	}

	hash, err := GetDDDataHash(context.Background(), db1)
	if err != nil {
		log.Errorf("cannot open dd-service database: %w", err)
		return
	}

	hash2, err := GetDDDataHash(context.Background(), db2)
	if err != nil {
		log.Errorf("cannot open dd-service database: %w", err)
		return
	}

	fmt.Println("hash1", hash)
	fmt.Println("hash2", hash2)
}

func GetDDDataHash(ctx context.Context, db *sqlx.DB) (hash string, err error) {
	r := []fingerprints{}
	err = db.Select(&r, "SELECT sha256_hash_of_art_image_file FROM image_hash_to_image_fingerprint_table where txid_timestamp >0 order by sha256_hash_of_art_image_file asc")
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get image_hash_to_image_fingerprint_table")
	}

	c := []collections{}
	err = db.Select(&c, "SELECT collection_ticket_txid FROM collections_table")
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

type fingerprints struct {
	Sha256HashOfArtImageFile string `db:"sha256_hash_of_art_image_file,omitempty"`
}
type collections struct {
	CollectionTicketTXID string `db:"collection_ticket_txid"`
}
