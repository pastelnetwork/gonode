package rqlite

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
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (s *testSuite) SetupSuite() {
	// new the rqlite config
	s.config = NewConfig()
	// new the rqlite service
	s.service = NewService(s.config)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// start the rqlite server
		if err := s.service.Run(ctx); err != nil {
			panic(err)
		}
	}()
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
	// create the table

}
