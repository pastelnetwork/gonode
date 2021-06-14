package kademlia

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

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
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
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
	res, err := ts.main.Store(ts.ctx, ts.Value)
	if err != nil {
		ts.T().Fatalf("dht store: %v", err)
	}

	key := base58.Decode(res)
	// retrieve the value for the key from storage
	actual := ts.main.store.Retrieve(ts.ctx, key)
	ts.Equal(ts.Value, actual)
}
