package kademlia

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/store/sqlite"

	"github.com/btcsuite/btcutil/base58"
	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/mem"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/suite"
	"go.uber.org/ratelimit"
)

var TestPri [144]byte
var TestPub [56]byte

func init() {
	curve := ed448.NewCurve()
	// generate the private and public key for client
	TestPri, TestPub, _ = curve.GenerateKeys()
}

// FakePastelClient - fake of pastel client to do sign/verify
type FakePastelClient struct {
	*pastelMock.Client
	curve ed448.Curve
	pri   [144]byte
	pub   [56]byte
}

// Sign
func (c *FakePastelClient) Sign(_ context.Context, data []byte, _, _ string, _ string) ([]byte, error) {
	signature, ok := c.curve.Sign(c.pri, data)
	if !ok {
		return nil, errors.New("sign failed")
	}
	signatureStr := base64.StdEncoding.EncodeToString(signature[:])
	return []byte(signatureStr), nil
}

// Verify
func (c *FakePastelClient) Verify(_ context.Context, data []byte, signature, _ string, _ string) (ok bool, err error) {
	signatureData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.Errorf("decode failed %w", err)
	}
	var copiedSignature [112]byte
	copy(copiedSignature[:], signatureData)

	ok = c.curve.Verify(copiedSignature, data, c.pub)
	return ok, nil
}

// Verify
func (c *FakePastelClient) MasterNodesExtra(_ context.Context) (pastel.MasterNodes, error) {
	return pastel.MasterNodes{
		pastel.MasterNode{
			ExtAddress: "127.0.0.1:1000",
			ExtKey:     "",
		},
	}, nil
}

func TestSuite(t *testing.T) {
	ts := new(testSuite)
	ts.t = t
	suite.Run(t, ts)
}

type testSuite struct {
	suite.Suite

	t       *testing.T
	Key     string
	Value   []byte
	IP      string
	DataDir string

	bootstrapPort  int
	bootstrapNodes []*Node

	dbStore  Store
	memStore Store
	main     *DHT
	ctx      context.Context
	cancel   context.CancelFunc
}

func (ts *testSuite) SetupSuite() {
	// init the log level
	if os.Getenv("LEVEL") == "debug" {
		log.SetP2PLogLevelName("debug")
	} else if os.Getenv("LEVEL") == "info" {
		log.SetP2PLogLevelName("info")
	} else {
		log.SetP2PLogLevelName("warn")
	}

	ts.IP = "127.0.0.1"
	ts.bootstrapPort = 6000
	ts.bootstrapNodes = []*Node{
		{
			IP:   ts.IP,
			Port: ts.bootstrapPort,
		},
	}
	// init the data directory for each test
	workDir, err := ioutil.TempDir("", "p2p-test-*")
	if err != nil {
		ts.T().Fatalf("ioutil tempdir: %v", err)
	}

	ts.ctx, ts.cancel = context.WithCancel(context.Background())
	ts.ctx = log.ContextWithPrefix(ts.ctx, "p2p-test")

	// init the badger store
	defaultReplicateInterval := time.Second * 3600
	defaultRepublishInterval := time.Second * 3600 * 24

	dbStore, err := sqlite.NewStore(ts.ctx, filepath.Join(workDir, "p2p"), defaultReplicateInterval, defaultRepublishInterval)
	if err != nil {
		ts.T().Fatalf("new sqlite store: %v", err)
	}
	ts.dbStore = dbStore
	ts.memStore = mem.NewStore()

	// init the main dht
	dht, err := ts.newDHTNodeWithDBStore(ts.ctx, ts.bootstrapPort, nil, nil)
	if err != nil {
		ts.T().Fatalf("new main dht: %v", err)
	}
	// reset the rate limiter
	dht.network.limiter = ratelimit.NewUnlimited()

	// start the main dht
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start main dht: %v", err)
	}
	ts.main = dht

	// init key and value for one data
	ts.Value = make([]byte, 50*1024)
	rand.Read(ts.Value)
	ts.Key = base58.Encode(dht.hashKey(ts.Value))
}

// run before each test in the suite
func (ts *testSuite) SetupTest() {
	// make sure the store is empty
	keys, _ := ts.main.store.Count(ts.ctx)
	ts.Zero(keys)
}

// reset the hashtable
func (ts *testSuite) resetHashtable() {
	ts.main.ht.refreshers = make([]time.Time, B)
	ts.main.ht.routeTable = make([][]*Node, B)
}

// run after each test in the suite
func (ts *testSuite) TearDownTest() {
	// reset the hashtable
	ts.resetHashtable()

	// reset the store
	_ = ts.main.store.DeleteAll(ts.ctx)
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	if ts.main != nil {
		ts.main.Stop(ts.ctx)
	}
	if ts.dbStore != nil {
		ts.dbStore.Close(ts.ctx)
	}
}

func (ts *testSuite) newDHTNodeWithMemStore(ctx context.Context, port int, nodes []*Node, id []byte) (*DHT, error) {
	options := &Options{
		IP:   ts.IP,
		Port: port,
	}
	if len(id) > 0 {
		options.ID = id
	}
	if len(nodes) > 0 {
		options.BootstrapNodes = nodes
	}

	pnodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		pnodes = append(pnodes, pastel.MasterNode{})
	}
	pnodes = append(pnodes, pastel.MasterNode{ExtKey: "A"})

	pastelClientMock := pastelMock.NewMockClient(ts.t)
	pastelClientMock.ListenOnMasterNodesTop(pnodes, nil)

	secInfo := &alts.SecInfo{}

	// fakePastelClient := &FakePastelClient{
	// 	curve: ed448.NewCurve(),
	// 	pri:   TestPri,
	// 	pub:   TestPub,
	// }
	// transportCredentials := credentials.NewClientCreds(fakePastelClient, secInfo)

	dht, err := NewDHT(ctx, ts.memStore, pastelClientMock, secInfo, options)
	if err != nil {
		return nil, errors.Errorf("new dht: %w", err)
	}

	return dht, nil
}

func (ts *testSuite) newDHTNodeWithDBStore(ctx context.Context, port int, nodes []*Node, id []byte) (*DHT, error) {
	options := &Options{
		IP:       ts.IP,
		Port:     port,
		PeerAuth: true, // Enable peer authentication
	}
	if len(id) > 0 {
		options.ID = id
	}
	if len(nodes) > 0 {
		options.BootstrapNodes = nodes
	}

	pnodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		pnodes = append(pnodes, pastel.MasterNode{})
	}
	pnodes = append(pnodes, pastel.MasterNode{ExtKey: "A"})

	pastelClientMock := pastelMock.NewMockClient(ts.t)
	pastelClientMock.ListenOnMasterNodesTop(pnodes, nil)

	secInfo := &alts.SecInfo{
		PastelID: "",
	}

	fakePastelClient := &FakePastelClient{
		Client: pastelClientMock,
		curve:  ed448.NewCurve(),
		pri:    TestPri,
		pub:    TestPub,
	}

	dht, err := NewDHT(ctx, ts.dbStore, fakePastelClient, secInfo, options)
	if err != nil {
		return nil, errors.Errorf("new dht: %w", err)
	}
	dht.network.limiter = ratelimit.NewUnlimited()

	return dht, nil
}

func (ts *testSuite) TestStartDHTNode() {
	// new a dht node
	dht, err := ts.newDHTNodeWithDBStore(ts.ctx, 7000, ts.bootstrapNodes, nil)
	if err != nil {
		ts.T().Fatalf("start a dht: %v", err)
	}

	// do the bootstrap
	if err := dht.Bootstrap(ts.ctx, ""); err != nil {
		ts.T().Fatalf("do bootstrap: %v", err)
	}

	// start the dht node
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start dht node: %v", err)
	}
	defer dht.Stop(ts.ctx)
}

func (ts *testSuite) TestRetrieveWithNil() {
	value, err := ts.main.Retrieve(ts.ctx, string(ts.Key))
	if err != nil {
		ts.T().Fatalf("dht retrive: %v", err)
	}
	ts.Nil(value)
}

func (ts *testSuite) TestStoreAndRetrieveWithMemStore() {
	// new a dht node
	dht, err := ts.newDHTNodeWithMemStore(ts.ctx, 6001, ts.bootstrapNodes, nil)
	if err != nil {
		ts.T().Fatalf("start a dht: %v", err)
	}

	// do the bootstrap
	if err := dht.Bootstrap(ts.ctx, ""); err != nil {
		ts.T().Fatalf("do bootstrap: %v", err)
	}

	// start the dht node
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start dht node: %v", err)
	}
	defer dht.Stop(ts.ctx)

	key, err := dht.Store(ts.ctx, ts.Value, "test-data-type")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	value, err := dht.Retrieve(ts.ctx, key)
	if err != nil {
		ts.T().Fatalf("dht retrieve: %v", err)
	}
	ts.Equal(ts.Value, value)
}

func (ts *testSuite) TestStoreAndRetrieve() {
	key, err := ts.main.Store(ts.ctx, ts.Value, "")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	value, err := ts.main.Retrieve(ts.ctx, key)
	if err != nil {
		ts.T().Fatalf("dht retrieve: %v", err)
	}
	ts.Equal(ts.Value, value)
}

func (ts *testSuite) TestStoreWithMain() {
	encodedKey, err := ts.main.Store(ts.ctx, ts.Value, "")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	// verify the value for dht store
	ts.verifyValue(base58.Decode(encodedKey), ts.Value, ts.main)
}

func (ts *testSuite) TestStoreWithTwoNodes() {
	// new a dht node
	dht, err := ts.newDHTNodeWithDBStore(ts.ctx, 6001, ts.bootstrapNodes, nil)
	if err != nil {
		ts.T().Fatalf("start a dht: %v", err)
	}

	// do the bootstrap
	if err := dht.Bootstrap(ts.ctx, ""); err != nil {
		ts.T().Fatalf("do bootstrap: %v", err)
	}

	// start the dht node
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start dht node: %v", err)
	}
	defer dht.Stop(ts.ctx)

	// store the key and value by main node
	encodedKey, err := ts.main.Store(ts.ctx, ts.Value, "")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}
	ts.Equal(ts.Key, encodedKey)

	// verify the value for dht store
	ts.verifyValue(base58.Decode(encodedKey), ts.Value, ts.main, dht)
}

// verify the value stored for distributed hash table
func (ts *testSuite) verifyValue(key []byte, value []byte, dhts ...*DHT) {
	for _, dht := range dhts {
		data, err := dht.store.Retrieve(ts.ctx, key)
		if err != nil {
			ts.T().Fatalf("store retrieve: %v", err)
		}
		ts.Equal(value, data)
	}
}

// start n nodes
func (ts *testSuite) startNodes(start, end int) ([]*DHT, error) {
	dhts := []*DHT{}

	for i := start; i <= end; i++ {
		// new a dht node
		dht, err := ts.newDHTNodeWithDBStore(ts.ctx, ts.bootstrapPort+i+1, ts.bootstrapNodes, nil)
		if err != nil {
			return nil, errors.Errorf("start a dht: %w", err)
		}

		// do the bootstrap
		if err := dht.Bootstrap(ts.ctx, ""); err != nil {
			return nil, errors.Errorf("do bootstrap: %w", err)
		}

		// start the dht node
		if err := dht.Start(ts.ctx); err != nil {
			return nil, errors.Errorf("start dht node: %w", err)
		}
		dhts = append(dhts, dht)
	}
	return dhts, nil
}

func (ts *testSuite) TestStoreWith10Nodes() {
	start, end := 1, 10
	// start 10 nodes
	dhts, err := ts.startNodes(start, end)
	if err != nil {
		ts.T().Fatalf("start nodes: %v", err)
	}
	for _, dht := range dhts {
		defer dht.Stop(ts.ctx)
	}
	dhts = append(dhts, ts.main)

	// init key and value for one data
	value := make([]byte, 50*1024)
	rand.Read(value)
	key := base58.Encode(ts.main.hashKey(value))

	// store the key and value by main node
	encodedKey, err := ts.main.Store(ts.ctx, value, "")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}
	ts.Equal(key, encodedKey)

	// verify the value for dht store
	ts.verifyValue(base58.Decode(encodedKey), value, dhts...)
}

func (ts *testSuite) TestIterativeFindNode() {
	start, end := 11, 20
	// start the nodes
	dhts, err := ts.startNodes(start, end)
	if err != nil {
		ts.T().Fatalf("start nodes: %v", err)
	}
	for _, dht := range dhts {
		defer dht.Stop(ts.ctx)
	}

	if _, err := ts.main.iterate(ts.ctx, IterateFindNode, ts.main.ht.self.ID, nil); err != nil {
		ts.T().Fatalf("iterative find node: %v", err)
	}

	for _, dht := range dhts {
		ts.Equal(end-start+1, dht.ht.totalCount())
	}
}

func (ts *testSuite) TestIterativeFindValue() {
	start, end := 21, 22
	// start the nodes
	dhts, err := ts.startNodes(start, end)
	if err != nil {
		ts.T().Fatalf("start nodes: %v", err)
	}
	for _, dht := range dhts {
		defer dht.Stop(ts.ctx)
	}

	// store the key/value
	key, err := ts.main.Store(ts.ctx, ts.Value, "")
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	value, err := ts.main.iterate(ts.ctx, IterateFindValue, base58.Decode(key), nil)
	if err != nil {
		ts.T().Fatalf("iterative find value: %v", err)
	}
	ts.Equal(ts.Value, value)
}

func (ts *testSuite) TestAddNodeForFull() {
	// short the ping time
	defaultPingTime = time.Millisecond * 200

	full := false

	for i := 0; i < 100; i++ {
		id, _ := newRandomID()
		node := &Node{
			ID:   id,
			IP:   ts.IP,
			Port: ts.bootstrapPort + i + 1,
		}
		if removed := ts.main.addNode(context.Background(), node); removed != nil {
			ts.Equal(false, ts.main.ht.hasNode(removed.ID))
			full = true
			break
		}
	}
	ts.Equal(true, full)

	full = false
	for _, bucket := range ts.main.ht.routeTable {
		if len(bucket) == K {
			full = true
		}
	}
	ts.Equal(true, full)
}

func (ts *testSuite) TestAddNodeForRefresh() {
	number := 20
	for i := 0; i < number; i++ {
		id, _ := newRandomID()
		node := &Node{
			ID:   id,
			IP:   ts.IP,
			Port: ts.bootstrapPort + i + 1,
		}
		ts.main.addNode(context.Background(), node)
	}
	ts.Equal(number, ts.main.ht.totalCount())

	// find the bucket which has the maximum number of nodes
	index := 0
	for i, bucket := range ts.main.ht.routeTable {
		if len(bucket) > 0 {
			index = i
		}
	}
	bucket := ts.main.ht.routeTable[index]
	count := len(bucket)
	first := bucket[0]

	// refresh the hash table with the first node
	ts.main.addNode(context.Background(), bucket[0])
	bucket = ts.main.ht.routeTable[index]
	ts.Equal(count, len(bucket))
	ts.Equal(first.ID, bucket[len(bucket)-1].ID)
}

func (ts *testSuite) TestAddNodeForAppend() {
	number := 20
	for i := 0; i < number; i++ {
		id, err := newRandomID()
		if err != nil {
			ts.T().Fatalf("new random id: %v", err)
		}
		node := &Node{
			ID:   id,
			IP:   ts.IP,
			Port: ts.bootstrapPort + i + 1,
		}
		ts.main.addNode(context.Background(), node)
	}
	ts.Equal(number, ts.main.ht.totalCount())

	// calculate the bucket index for a new random id
	id, err := newRandomID()
	if err != nil {
		ts.T().Fatalf("new random id: %v", err)
	}
	index := ts.main.ht.bucketIndex(ts.main.ht.self.ID, id)

	// add the node to hash table
	node := &Node{
		ID:   id,
		IP:   ts.IP,
		Port: ts.bootstrapPort + number + 1,
	}
	ts.main.addNode(context.Background(), node)

	bucket := ts.main.ht.routeTable[index]
	ts.Equal(id, bucket[len(bucket)-1].ID)
}

func (ts *testSuite) TestHashKey() {
	key := ts.main.hashKey(ts.Value)
	ts.Equal(B/8, len(key))
}

func (ts *testSuite) TestRateLimiter() {
	start, end := 23, 32
	// start 10 nodes
	dhts, err := ts.startNodes(start, end)
	if err != nil {
		ts.T().Fatalf("start nodes: %v", err)
	}
	for _, dht := range dhts {
		defer dht.Stop(ts.ctx)
	}
}
