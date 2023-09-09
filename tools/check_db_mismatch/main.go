package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func fetchTxIDsFromDB(db *sql.DB) (map[string]bool, error) {
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
	db1, err := sql.Open("sqlite3", "./hermes-154.12.235.16.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db1.Close()

	db2, err := sql.Open("sqlite3", "./hermes-154.12.235.41.sqlite")
	if err != nil {
		log.Fatal(err)
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

	fmt.Println("TxIDs in hermes-154.12.235.16.sqlite but not in hermes-154.12.235.41.sqlite:")
	for txID := range txIDs1 {
		if _, exists := txIDs2[txID]; !exists {
			fmt.Println(txID)
		}
	}

	fmt.Println("\nTxIDs in hermes-154.12.235.41.sqlite but not in hermes-154.12.235.16.sqlite:")
	for txID := range txIDs2 {
		if _, exists := txIDs1[txID]; !exists {
			fmt.Println(txID)
		}
	}
}
