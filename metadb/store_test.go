package metadb

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAndQuery(t *testing.T) {
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

	// create the table
	creation := "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)"
	result, err := s.Write(ctx, creation)
	assert.Nil(t, err, "execute creation statement")
	assert.Empty(t, result.Error, "result after execution")

	// insert the value to table
	insertion := "insert into foo(id, name) values(1,'alon')"
	result, err = s.Write(ctx, insertion)
	assert.Nil(t, err, "execute insertion statement")
	assert.Empty(t, result.Error, "result after execution")

	// insert the value to table
	insertion = "insert into foo(id, name) values(2,'belayu')"
	result, err = s.Write(ctx, insertion)
	assert.Nil(t, err, "execute insertion statement")
	assert.Empty(t, result.Error, "result after execution")

	// select a row from table with week level
	selection := "select * from foo where id = 1"
	rows, err := s.Query(ctx, selection, "weak")
	assert.Nil(t, err, "query statements")

	var id int64
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		assert.Nil(t, err, "rows scan")
		assert.Equal(t, int64(1), id)
		assert.Equal(t, "alon", name)
	}

	cancel()
	wg.Wait()

	// remove the data directy
	os.RemoveAll(s.config.DataDir)
}
