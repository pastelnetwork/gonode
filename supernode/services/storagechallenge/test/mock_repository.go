package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	GetListOfMasternodeMethod                             = "GetListOfMasternode"
	GetSymbolFileByKeyMethod                              = "GetSymbolFileByKey"
	ListSymbolFileKeysFromNFTTicketMethod                 = "ListSymbolFileKeysFromNFTTicket"
	GetNClosestMasternodeIDsToComparisionStringMethod     = "GetNClosestMasternodeIDsToComparisionString"
	GetNClosestMasternodesToAGivenFileUsingKademliaMethod = "GetNClosestMasternodesToAGivenFileUsingKademlia"
	GetNClosestFileHashesToAGivenComparisonStringMethod   = "GetNClosestFileHashesToAGivenComparisonString"
	SaveChallengMessageStateMethod                        = "SaveChallengMessageState"
)

type Repo struct {
	t *testing.T
	*mocks.Repository
}

func (r *Repo) ListenOnGetListOfMasternode(mnList []string, err error) *Repo {
	r.On(GetListOfMasternodeMethod, mock.Anything).Return(mnList, err)
	return r
}

func (r *Repo) ListenOnGetSymbolFileByKey(value []byte, err error) *Repo {
	r.On(GetSymbolFileByKeyMethod, mock.Anything, mock.Anything, mock.Anything).Return(value, err)
	return r
}

func (r *Repo) ListenOnListSymbolFileKeysFromNFTTicket(listSymbolFileHashes []string, err error) *Repo {
	r.On(ListSymbolFileKeysFromNFTTicketMethod, mock.Anything).Return(listSymbolFileHashes, err)
	return r
}

func (r *Repo) ListenOnGetNClosestMasternodeIDsToComparisionString(closest []string) *Repo {
	r.On(GetNClosestMasternodeIDsToComparisionStringMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(closest)
	return r
}

func (r *Repo) ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore(closest []string) *Repo {
	r.On(GetNClosestMasternodeIDsToComparisionStringMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(closest)
	return r
}

func (r *Repo) ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia(closest []string) *Repo {
	r.On(GetNClosestMasternodesToAGivenFileUsingKademliaMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(closest)
	return r
}

func (r *Repo) ListenOnGetNClosestFileHashesToAGivenComparisonString(closest []string) *Repo {
	r.On(GetNClosestFileHashesToAGivenComparisonStringMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(closest)
	return r
}

func (r *Repo) ListenOnSaveChallengMessageState() *Repo {
	r.On(SaveChallengMessageStateMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return r
}

func (r *Repo) AssertGetListOfMasternodeCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetListOfMasternodeMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetListOfMasternodeMethod, expectedCalls)
	return r
}

func (r *Repo) AssertGetSymbolFileByKeyCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetSymbolFileByKeyMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetSymbolFileByKeyMethod, expectedCalls)
	return r
}

func (r *Repo) AssertListSymbolFileKeysFromNFTTicketCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, ListSymbolFileKeysFromNFTTicketMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, ListSymbolFileKeysFromNFTTicketMethod, expectedCalls)
	return r
}

func (r *Repo) AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetNClosestMasternodeIDsToComparisionStringMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetNClosestMasternodeIDsToComparisionStringMethod, expectedCalls)
	return r
}

func (r *Repo) AssertGetNClosestMasternodeIDsToComparisionStringCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetNClosestMasternodeIDsToComparisionStringMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetNClosestMasternodeIDsToComparisionStringMethod, expectedCalls)
	return r
}

func (r *Repo) AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetNClosestMasternodesToAGivenFileUsingKademliaMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetNClosestMasternodesToAGivenFileUsingKademliaMethod, expectedCalls)
	return r
}

func (r *Repo) AssertGetNClosestFileHashesToAGivenComparisonStringCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, GetNClosestFileHashesToAGivenComparisonStringMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, GetNClosestFileHashesToAGivenComparisonStringMethod, expectedCalls)
	return r
}

func (r *Repo) AssertSaveChallengMessageStateCall(expectedCalls int, arguments ...interface{}) *Repo {
	if expectedCalls > 0 {
		r.AssertCalled(r.t, SaveChallengMessageStateMethod, arguments...)
	}
	r.AssertNumberOfCalls(r.t, SaveChallengMessageStateMethod, expectedCalls)
	return r
}

func NewMockRepository(t *testing.T) *Repo {
	return &Repo{
		t:          t,
		Repository: &mocks.Repository{},
	}
}
