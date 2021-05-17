package kademlia

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Test_migrate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
}

func Test_store(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	if err := store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}
}
