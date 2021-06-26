package metadb

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/log"
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
	log.SetOutput(ioutil.Discard)

	workDir, err := ioutil.TempDir("", "metadb-*")
	assert.NoError(ts.T(), err)

	config := NewConfig()
	config.SetWorkDir(workDir)

	ts.s = &service{
		nodeID: "node",
		config: config,
		ready:  make(chan struct{}, 5),
	}
	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		// start the rqlite server
		ts.Nil(ts.s.Run(ts.ctx), "run service")
	}()
	<-ts.s.ready

	tableThumbnails := `
		CREATE TABLE IF NOT EXISTS thumbnails (
			key TEXT PRIMARY KEY,
			small BLOB,
			medium BLOB,
			large BLOB
		)
	`
	// create the thumbnails table
	if _, err := ts.s.Write(ts.ctx, tableThumbnails); err != nil {
		ts.T().Fatalf("write: %v", err)
	}

	if _, err := ts.s.Write(ts.ctx, "CREATE TABLE IF NOT EXISTS foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}
	if _, err := ts.s.Write(ts.ctx, "insert into foo(id, name) values(1,'1')"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}
	if _, err := ts.s.Write(ts.ctx, "insert into foo(id, name) values(2,'2')"); err != nil {
		ts.T().Fatalf("write: %v", err)
	}

	ts.NotNil(ts.s.db, "store is nil")
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.wg.Wait()

	// remove the data directy
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
			ts.Equal(tc.err, err)
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
			ts.Nil(err, "query statement")

			var id int64
			var name string
			for rows.Next() {
				ts.Equal(tc.err, rows.Scan(&id, &name))
				ts.Equal(tc.id, id)
				ts.Equal(tc.name, name)
			}
		})
	}
}

// Thumbnail stucture for rqlite table 'thumbnails'
type Thumbnail struct {
	Key    string `json:"key"`
	Small  []byte `json:"small"`
	Medium []byte `json:"medium"`
	Large  []byte `json:"large"`
}

func (ts *testSuite) randomCasesForThumbnails(start, end int) []*Thumbnail {
	small, medium, large := [10]byte{}, [15]byte{}, [20]byte{}

	thumbnails := []*Thumbnail{}
	for i := start; i <= end; i++ {
		rand.Read(small[:])
		rand.Read(medium[:])
		rand.Read(large[:])

		thumbnails = append(thumbnails, &Thumbnail{
			Key:    "key-" + strconv.Itoa(i),
			Small:  small[:],
			Medium: medium[:],
			Large:  large[:],
		})
	}
	return thumbnails
}

func (ts *testSuite) TestThumbnailTable() {
	for i, tc := range ts.randomCasesForThumbnails(5, 10) {
		tc := tc
		ts.T().Run(fmt.Sprintf("ThumbnailTestCase-%d", i), func(t *testing.T) {
			// insert the thumbnail to db
			insertion := "insert into thumbnails(key,small,medium,large) values(?,?,?,?)"
			if _, err := ts.s.Write(ts.ctx, insertion, tc.Key, tc.Small, tc.Medium, tc.Large); err != nil {
				ts.T().Fatalf("write: %v", err)
			}

			// select the thumbnail by key
			query := fmt.Sprintf("select * from thumbnails where key='%s'", tc.Key)
			rows, err := ts.s.Query(ts.ctx, query, "weak")
			ts.Nil(err, "query statement")
			ts.NotNil(rows)

			var result Thumbnail
			for rows.Next() {
				ts.Nil(rows.Scan(&result.Key, &result.Small, &result.Medium, &result.Large))
			}
			ts.Equal(tc.Key, result.Key)
			ts.Equal(tc.Small, result.Small)
			ts.Equal(tc.Medium, result.Medium)
			ts.Equal(tc.Large, result.Large)
		})
	}
}

func (ts *testSuite) TestAddThumbnailTwice() {
	for i, tc := range ts.randomCasesForThumbnails(1, 1) {
		tc := tc
		ts.T().Run(fmt.Sprintf("ThumbnailTestCase-%d", i), func(t *testing.T) {
			// insert the thumbnail to db
			insertion := "insert into thumbnails(key,small,medium,large) values(?,?,?,?)"
			if _, err := ts.s.Write(ts.ctx, insertion, tc.Key, tc.Small, tc.Medium, tc.Large); err != nil {
				ts.T().Fatalf("write: %v", err)
			}

			// insert the same thumbnail again
			_, err := ts.s.Write(ts.ctx, insertion, tc.Key, tc.Small, tc.Medium, tc.Large)
			ts.NotNil(err)
		})
	}
}
