package storagechallenge

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func (suite *testSuite) SetupTest() {
}

func (suite *testSuite) TearDownTest() {
}
func (suite *testSuite) TestFullChallengeSucceeded() {
	// implement the test here
}

func (suite *testSuite) TestDishonestNodeByFailedOrTimeout() {
	// implement the test here
}

func (suite *testSuite) TestIncrementalRQSymbolFileKey() {
	// implement the test here
}

func (suite *testSuite) TestIncrementalMasternode() {
	// implement the test here
}

func TestStorageChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
