package healthcheck

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/mr-tron/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb"
	"golang.org/x/crypto/sha3"
)

const (
	// B - the size in bits of the keys used to identify nodes and store and
	// retrieve data; in basic Kademlia this is 256, the length of a SHA3-256
	B = 256
)

//  generate hash value for the data
func sha3256Hash(data []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func rqliteKeyGen(data []byte) (string, error) {
	hash, err := sha3256Hash(data)
	if err != nil {
		return "", err
	}

	key := base58.Encode(hash)
	return key, nil
}

// keyValueResult represents a item of store
type keyValueResult struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type RqliteKVStore struct {
	metaDb metadb.MetaDB
}

func NewRqliteKVStore(metaDb metadb.MetaDB) *RqliteKVStore {
	return &RqliteKVStore{
		metaDb: metaDb,
	}
}

// Retrieve retrieve data from the kademlia network by key, return nil if not found
// - key is the base58 encoded identifier of the data
func (store *RqliteKVStore) Retrieve(ctx context.Context, key string) ([]byte, error) {
	queryString := `SELECT * FROM test_keystore WHERE key = "` + key + `"`

	// do query
	queryResult, err := store.metaDb.Query(ctx, queryString, "nono")
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return nil, errors.New("empty result")
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()
	resultMap, err := queryResult.Map()
	if err != nil {
		return nil, errors.Errorf("query map: %w", err)
	}

	var dbResult keyValueResult
	if err := mapstructure.Decode(resultMap, &dbResult); err != nil {
		return nil, errors.Errorf("mapstructure decode: %w", err)
	}

	// decode value
	value, err := hex.DecodeString(dbResult.Value)
	if err != nil {
		return nil, errors.Errorf("hex decode: %w", err)
	}

	return value, nil

}

// Store store data to the network, which will trigger the iterative store message
func (store *RqliteKVStore) Store(ctx context.Context, data []byte) (string, error) {
	key, err := rqliteKeyGen(data)
	if err != nil {
		return "", errors.Errorf("gen key: %v", err)
	}

	// convert string to hex format - simple technique to prevent SQL security issue
	hexValue := hex.EncodeToString(data)
	queryString := `INSERT INTO test_keystore(key, value) VALUES("` + key + `", "` + hexValue + `")`

	// do query
	_, err = store.metaDb.Query(ctx, queryString, "nono")
	if err != nil {
		return "", errors.Errorf("querying db: %w", err)
	}

	return key, nil
}

// Delete a key, value
func (store *RqliteKVStore) Delete(ctx context.Context, key string) error {
	// verify key
	decoded, err := base58.Decode(key)
	if err != nil || len(decoded) != B/8 {
		return fmt.Errorf("invalid key: %v", key)
	}

	queryString := `DELETE FROM test_keystore WHERE key = "` + key + `"`

	// do query
	_, err = store.metaDb.Query(ctx, queryString, "nono")
	if err != nil {
		return errors.Errorf("querying db: %w", err)
	}

	return nil
}
