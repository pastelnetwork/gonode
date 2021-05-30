package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func setupDB() *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")
	Migrate(ctx, db)
	return db
}

func TestMigrate(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "running migrations",
			args:      args{db},
			assertion: assert.NoError,
		},
		{
			name:      "running migrations twice",
			args:      args{db},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, Migrate(ctx, tt.args.db))
		})
	}
}

func TestStore(t *testing.T) {
	db := setupDB()

	type args struct {
		db          *sql.DB
		data        []byte
		replication time.Time
		expiration  time.Time
	}
	tests := []struct {
		name      string
		args      args
		want      int64
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "storing one key",
			args:      args{db, []byte("test data"), time.Now(), time.Now()},
			want:      int64(1),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Store(ctx, tt.args.db, tt.args.data, tt.args.replication, tt.args.expiration)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRetrieve(t *testing.T) {
	db := setupDB()
	_, err := Store(ctx, db, []byte("test data"), time.Now(), time.Now())
	assert.NoError(t, err)

	type args struct {
		db  *sql.DB
		key []byte
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "good key",
			args:      args{db, crypto.GetKey([]byte("test data"))},
			want:      []byte("test data"),
			assertion: assert.NoError,
		},
		{
			name:      "bad key",
			args:      args{db, []byte("12312312313123")},
			want:      []byte(nil),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := Retrieve(ctx, tt.args.db, tt.args.key)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpireKeys(t *testing.T) {
	db := setupDB()
	_, err := Store(ctx, db, []byte("test data 1"), time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = Store(ctx, db, []byte("test data 2"), time.Now(), time.Now())
	assert.NoError(t, err)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "expire test data",
			args:      args{db},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, ExpireKeys(ctx, tt.args.db))
		})
	}
}

func TestGetAllKeysForReplication(t *testing.T) {
	db := setupDB()
	_, err := Store(ctx, db, []byte("test data 1"), time.Now(), time.Now())
	assert.NoError(t, err)
	_, err = Store(ctx, db, []byte("test data 2"), time.Now(), time.Now())
	assert.NoError(t, err)
	_, err = Store(ctx, db, []byte("test data 2"), time.Now().Add(1*time.Hour), time.Now().Add(1*time.Hour))
	assert.NoError(t, err)
	// should not return ^

	var want [][]byte
	want = append(want, crypto.GetKey([]byte("test data 1")))
	want = append(want, crypto.GetKey([]byte("test data 2")))

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		want      [][]byte
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "get all keys",
			args:      args{db},
			want:      want,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllKeysForReplication(ctx, tt.args.db)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemove(t *testing.T) {
	db := setupDB()
	_, err := Store(ctx, db, []byte("test data 1"), time.Now(), time.Now())
	assert.NoError(t, err)

	type args struct {
		db  *sql.DB
		key []byte
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "remove a key",
			args:      args{db, crypto.GetKey([]byte("test data 1"))},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, Remove(ctx, tt.args.db, tt.args.key))
		})
	}
}
