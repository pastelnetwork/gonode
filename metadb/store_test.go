package metadb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite

	s      *service
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (ts *testSuite) SetupSuite() {
	workDir, err := ioutil.TempDir("", "metadb-*")
	assert.NoError(ts.T(), err)

	config := NewConfig()
	config.SetWorkDir(workDir)

	ts.s = &service{
		nodeID: "uuid",
		config: config,
		ready:  make(chan struct{}, 1),
	}
	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		// start the rqlite server
		ts.Nil(ts.s.Run(ts.ctx), "run service")
	}()
	<-ts.s.ready

	ts.NotNil(ts.s.db, "store is nil")

	if _, err := ts.s.Write(ts.ctx, "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}
	if _, err := ts.s.Write(ts.ctx, "insert into foo(id, name) values(1,'1')"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}
	if _, err := ts.s.Write(ts.ctx, "insert into foo(id, name) values(2,'2')"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.wg.Wait()

	// remove the data directly
	os.RemoveAll(filepath.Join(ts.s.config.DataDir))
}

func (ts *testSuite) TestWrite() {
	writeTestCases := []struct {
		stmt string
		err  error
	}{
		{
			stmt: "insert into foo(id, name) values(3,'3')",
		},
		{
			stmt: "insert into foo(id, name) values(4,'4')",
		},
	}
	for i, tc := range writeTestCases {
		tc := tc
		ts.T().Run(fmt.Sprintf("WriteTestCase-%d", i), func(t *testing.T) {
			_, err := ts.s.Write(ts.ctx, tc.stmt)
			assert.Equal(t, tc.err, err)
		})
	}
}

func (ts *testSuite) TestQuery() {
	queryTestCases := []struct {
		stmt string
		id   int64
		name string
		err  error
	}{
		{
			stmt: "select * from foo where id = 1",
			id:   1,
			name: "1",
		},
		{
			stmt: "select * from foo where id = 2",
			id:   2,
			name: "2",
		},
	}
	for i, tc := range queryTestCases {
		tc := tc
		ts.T().Run(fmt.Sprintf("QueryTestCase-%d", i), func(t *testing.T) {
			rows, err := ts.s.Query(ts.ctx, tc.stmt, "weak")
			assert.Nil(t, err, "query statement")

			var id int64
			var name string
			for rows.Next() {
				assert.Equal(t, tc.err, rows.Scan(&id, &name))
				assert.Equal(t, tc.id, id)
				assert.Equal(t, tc.name, name)
			}
		})
	}
}
