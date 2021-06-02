package rqlite

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/log"
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

	// disable the logger
	log.DefaultLogger.SetOutput(ioutil.Discard)
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

func (s *testSuite) TestQuery() {
}

func (s *testSuite) TestExecute() {
	s.NotNil(s.service.db, "store is nil")

	// create the table
	creation := "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)"
	// insert the value to table
	insertion := "insert into foo(id, name) values(1,'alon');insert into foo(id, name) values(2,'belayu')"

	stmts := []string{creation, insertion}
	// create table and insert the values
	results, err := s.service.Execute(s.ctx, stmts, true, true)
	s.Nil(err, "execute statements")
	for _, result := range results {
		s.Empty(result.Error, "execute result")
	}

	// select a row from table with week level
	selection := "select * from foo where id = 1"
	stmts = []string{selection}
	rows, err := s.service.Query(s.ctx, stmts, true, true, "weak")
	s.Nil(err, "query statements")
	s.Equal(1, len(rows))
	// for _, row := range rows {
	// }
}
