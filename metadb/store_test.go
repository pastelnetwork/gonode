package metadb

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAndQuery(t *testing.T) {
	t.Parallel()

	// start the rqlite node
	s := &service{
		config: NewConfig(),
		ready:  make(chan struct{}, 1),
	}
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// start the rqlite server
		assert.Nil(t, s.Run(ctx), "run service")
	}()
	// wait until the rqlite node is ready
	<-s.ready

	assert.NotNil(t, s.db, "store is nil")

	writeTestCases := []struct {
		stmt string
		err  error
	}{
		{
			stmt: "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)",
		},
		{
			stmt: "insert into foo(id, name) values(1,'alon')",
		},
		{
			stmt: "insert into foo(id, name) values(2,'belayu')",
		},
	}
	for i, tc := range writeTestCases {
		tc := tc
		t.Run(fmt.Sprintf("WriteTestCase-%d", i), func(t *testing.T) {
			_, err := s.Write(ctx, tc.stmt)
			assert.Equal(t, tc.err, err)
		})
	}

	queryTestCases := []struct {
		stmt string
		id   int64
		name string
		err  error
	}{
		{
			stmt: "select * from foo where id = 1",
			id:   1,
			name: "alon",
		},
		{
			stmt: "select * from foo where id = 2",
			id:   2,
			name: "belayu",
		},
	}
	for i, tc := range queryTestCases {
		tc := tc
		t.Run(fmt.Sprintf("QueryTestCase-%d", i), func(t *testing.T) {
			rows, err := s.Query(ctx, tc.stmt, "weak")
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

	cancel()
	wg.Wait()

	// remove the data directy
	os.RemoveAll(s.config.DataDir)
}
