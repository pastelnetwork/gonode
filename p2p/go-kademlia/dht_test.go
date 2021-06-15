package kademlia

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/jbenet/go-base58"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite

	Key   string
	Value []byte
	IP    string

	bootstrapPort  int
	bootstrapNodes []*Node

	main   *DHT
	ctx    context.Context
	cancel context.CancelFunc
}

func (ts *testSuite) SetupSuite() {
	// init the log formatter for logrus
	logrus.SetReportCaller(true)
	if os.Getenv("LEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if os.Getenv("LEVEL") == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, file, line, _ := runtime.Caller(10)
			return "", fmt.Sprintf(" %s:%v ", path.Base(file), line)
		},
	})

	ts.IP = "127.0.0.1"
	ts.bootstrapPort = 6000
	ts.bootstrapNodes = []*Node{
		{
			IP:   ts.IP,
			Port: ts.bootstrapPort,
		},
	}

	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	// init the main dht
	dht, err := ts.newDHTNode(ts.ctx, ts.bootstrapPort, nil, nil)
	if err != nil {
		ts.T().Fatalf("new main dht: %v", err)
	}
	// start the main dht
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start main dht: %v", err)
	}
	ts.main = dht

	// init key and value for one data
	ts.Value = []byte("hello world")
	ts.Key = base58.Encode(dht.hashKey(ts.Value))
}

// run before each test in the suite
func (ts *testSuite) SetupTest() {
	// make sure the store is empty
	keys := ts.main.store.Keys(ts.ctx)
	ts.Zero(len(keys))
}

// run after each test in the suite
func (ts *testSuite) TearDownTest() {
	// reset the hashtable
	ts.main.ht.reset()

	// reset the store
	for _, key := range ts.main.store.Keys(ts.ctx) {
		ts.main.store.Delete(ts.ctx, key)
	}
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	if ts.main != nil {
		ts.main.Stop(ts.ctx)
	}
}

func (ts *testSuite) newDHTNode(ctx context.Context, port int, nodes []*Node, id []byte) (*DHT, error) {
	store := NewMemStore()

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

	dht, err := NewDHT(store, options)
	if err != nil {
		return nil, fmt.Errorf("new dht: %v", err)
	}
	store.SetSelf(dht.ht.self)

	return dht, nil
}

func (ts *testSuite) TestStartDHTNode() {
	// new a dht node
	dht, err := ts.newDHTNode(ts.ctx, 6001, ts.bootstrapNodes, nil)
	if err != nil {
		ts.T().Fatalf("start a dht: %v", err)
	}

	// do the bootstrap
	if err := dht.Bootstrap(ts.ctx); err != nil {
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

func (ts *testSuite) TestStoreAndRetrieve() {
	key, err := ts.main.Store(ts.ctx, ts.Value)
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
	encodedKey, err := ts.main.Store(ts.ctx, ts.Value)
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	// verify the value for dht store
	ts.verifyValue(base58.Decode(encodedKey), ts.Value, ts.main)
}

func (ts *testSuite) TestStoreWithTwoNodes() {
	// new a dht node
	dht, err := ts.newDHTNode(ts.ctx, 6001, ts.bootstrapNodes, nil)
	if err != nil {
		ts.T().Fatalf("start a dht: %v", err)
	}

	// do the bootstrap
	if err := dht.Bootstrap(ts.ctx); err != nil {
		ts.T().Fatalf("do bootstrap: %v", err)
	}

	// start the dht node
	if err := dht.Start(ts.ctx); err != nil {
		ts.T().Fatalf("start dht node: %v", err)
	}
	defer dht.Stop(ts.ctx)

	// store the key and value by main node
	encodedKey, err := ts.main.Store(ts.ctx, ts.Value)
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
func (ts *testSuite) startNodes(n int) ([]*DHT, error) {
	dhts := []*DHT{}

	for i := 0; i < n; i++ {
		// new a dht node
		dht, err := ts.newDHTNode(ts.ctx, ts.bootstrapPort+i+1, ts.bootstrapNodes, nil)
		if err != nil {
			return nil, fmt.Errorf("start a dht: %v", err)
		}

		// do the bootstrap
		if err := dht.Bootstrap(ts.ctx); err != nil {
			return nil, fmt.Errorf("do bootstrap: %v", err)
		}

		// start the dht node
		if err := dht.Start(ts.ctx); err != nil {
			return nil, fmt.Errorf("start dht node: %v", err)
		}
		dhts = append(dhts, dht)
	}
	return dhts, nil
}

func (ts *testSuite) TestStoreWith10Nodes() {
	// start 10 nodes
	dhts, err := ts.startNodes(10)
	if err != nil {
		ts.T().Fatalf("start nodes: %v", err)
	}
	for _, dht := range dhts {
		defer dht.Stop(ts.ctx)
	}
	dhts = append(dhts, ts.main)

	// store the key and value by main node
	encodedKey, err := ts.main.Store(ts.ctx, ts.Value)
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}
	ts.Equal(ts.Key, encodedKey)

	// verify the value for dht store
	ts.verifyValue(base58.Decode(encodedKey), ts.Value, dhts...)
}

func (ts *testSuite) TestIterativeFindNode() {
	number := 2

	// start the nodes
	dhts, err := ts.startNodes(number)
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
		ts.Equal(number, dht.ht.totalCount())
	}
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
		if removed := ts.main.addNode(node); removed != nil {
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
		ts.main.addNode(node)
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
	first := bucket[0]

	// refresh the hash table with the first node
	ts.main.addNode(bucket[0])
	bucket = ts.main.ht.routeTable[index]
	last := bucket[len(bucket)-1]

	ts.Equal(first.ID, last.ID)
}
