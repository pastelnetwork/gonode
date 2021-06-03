package metadb

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite

	config  *Config
	service *Service
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (s *testSuite) SetupSuite() {
	// new the rqlite config
	s.config = NewConfig()
	// new the rqlite service
	s.service = NewService(s.config)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// start the rqlite server
		s.Nil(s.service.Run(s.ctx), "run service")
	}()
	// wait until the rqlite node is ready
	<-s.service.ready
}

func (s *testSuite) TearDownSuite() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()

	// remove the data directy
	os.RemoveAll(s.config.DataDir)
}

func (s *testSuite) TestWriteAndQuery() {
	s.NotNil(s.service.db, "store is nil")

	// create the table
	creation := "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)"
	result, err := s.service.Write(s.ctx, creation)
	s.Nil(err, "execute creation statement")
	s.Empty(result.Error, "result after execution")

	// insert the value to table
	insertion := "insert into foo(id, name) values(1,'alon')"
	result, err = s.service.Write(s.ctx, insertion)
	s.Nil(err, "execute insertion statement")
	s.Empty(result.Error, "result after execution")

	// insert the value to table
	insertion = "insert into foo(id, name) values(2,'belayu')"
	result, err = s.service.Write(s.ctx, insertion)
	s.Nil(err, "execute insertion statement")
	s.Empty(result.Error, "result after execution")

	// select a row from table with week level
	selection := "select * from foo where id = 1"
	rows, err := s.service.Query(s.ctx, selection, "weak")
	s.Nil(err, "query statements")

	var id int64
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		s.Nil(err, "rows scan")
		s.Equal(int64(1), id)
		s.Equal("alon", name)
	}
}
