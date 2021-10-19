package healthcheck

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mr-tron/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/p2p"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
)

const (
	defaultInputMaxLength  = 64
	defaultExpiresDuration = 15 * time.Minute
)

// StatsMngr is interface of StatsManger, return stats of system
type StatsMngr interface {
	// Stats returns stats of system
	Stats(ctx context.Context) (map[string]interface{}, error)
}

// P2PTracker is interface of lifetime tracker that will periodically check to remove expired temporary p2p's (key, value)
type P2PTracker interface {
	// Track monitor to delete temporary key
	Track(key string, expires time.Time)
}

// HealthCheck represents grpc service for supernode healthcheck
type HealthCheck struct {
	pb.UnimplementedHealthCheckServer
	StatsMngr
	p2pService p2p.P2P
	metaDb     metadb.MetaDB
	p2pTracker P2PTracker
}

// Status will send a message to and get back a reply from supernode
func (service *HealthCheck) Status(ctx context.Context, _ *pb.StatusRequest) (*pb.StatusReply, error) {
	stats, err := service.Stats(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to Stats(): %w", err)
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		return nil, errors.Errorf("failed to Marshal(): %w", err)
	}

	// echos received message
	return &pb.StatusReply{StatusInJson: string(jsonData)}, nil
}

// P2PStore store a data, and return given key
func (service *HealthCheck) P2PStore(ctx context.Context, in *pb.P2PStoreRequest) (*pb.P2PStoreReply, error) {
	value := in.GetValue()
	if len(value) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}
	key, err := service.p2pService.Store(ctx, []byte(value))

	if err != nil {
		return nil, err
	}

	// track to delete key soon
	service.p2pTracker.Track(key, time.Now().Add(defaultExpiresDuration))

	return &pb.P2PStoreReply{Key: key}, nil
}

// P2PRetrieve get data of given key
func (service *HealthCheck) P2PRetrieve(ctx context.Context, in *pb.P2PRetrieveRequest) (*pb.P2PRetrieveReply, error) {
	key := in.GetKey()
	value, err := service.p2pService.Retrieve(ctx, key)

	if err != nil {
		return &pb.P2PRetrieveReply{Value: ""}, err
	}
	return &pb.P2PRetrieveReply{Value: string(value)}, nil
}

//  generate hash value for the data
func sha3256Hash(data []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func rqliteKeyGen(msg string) (string, error) {
	hash, err := sha3256Hash([]byte(msg))
	if err != nil {
		return "", err
	}

	key := base58.Encode(hash)
	return key, nil
}

// RqliteStore do store a value - and return a key - is hash of string
func (service *HealthCheck) RqliteStore(ctx context.Context, in *pb.RqliteStoreRequest) (*pb.RqliteStoreReply, error) {
	if len(in.GetValue()) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}

	value := in.GetValue()
	key, err := rqliteKeyGen(value)
	if err != nil {
		return nil, errors.Errorf("gen key: %v", err)
	}

	// convert string to hex format - simple technique to prevent SQL security issue
	hexValue := hex.EncodeToString([]byte(value))
	queryString := `INSERT INTO test_keystore(key, value) VALUES("` + key + `", "` + hexValue + `")`

	// do query
	_, err = service.metaDb.Query(ctx, queryString, "nono")
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}

	return &pb.RqliteStoreReply{Key: key}, nil
}

// keyValueResult represents a item of store
type keyValueResult struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

// RqliteRetrieve get data of given key
func (service *HealthCheck) RqliteRetrieve(ctx context.Context, in *pb.RqliteRetrieveRequest) (*pb.RqliteRetrieveReply, error) {
	if len(in.GetKey()) > defaultInputMaxLength {
		return nil, errors.New("length too big (>32bytes")
	}

	key := in.GetKey()
	queryString := `SELECT * FROM test_keystore WHERE key = "` + key + `"`

	// do query
	queryResult, err := service.metaDb.Query(ctx, queryString, "nono")
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

	return &pb.RqliteRetrieveReply{Value: string(value)}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck(mngr StatsMngr, p2pService p2p.P2P, metaDb metadb.MetaDB, p2pTracker P2PTracker) *HealthCheck {
	return &HealthCheck{
		StatsMngr:  mngr,
		p2pService: p2pService,
		metaDb:     metaDb,
		p2pTracker: p2pTracker,
	}
}
