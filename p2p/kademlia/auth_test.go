package kademlia

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

// TestAuthenticator_isPeerInMasterNodes verify isPeerInMasterNodes()
func TestAuthenticator_isPeerInMasterNodes(t *testing.T) {
	type args struct {
		testInitLastMasterNodesRefresh time.Time
		testInitMasterNodesExtra       pastel.MasterNodes
		testRefreshMasterNodesExtra    pastel.MasterNodes
		testRefreshMasterNodesExtraErr error
		testPeerAuthInfo               *auth.PeerAuthInfo
	}

	masterNodes1 := pastel.MasterNode{
		ExtAddress: "192.168.1.1:1000",
		ExtKey:     "Pastel1",
	}
	masterNodes2 := pastel.MasterNode{
		ExtAddress: "192.168.1.2:1000",
		ExtKey:     "Pastel2",
	}

	testRefreshDuration := 5 * time.Minute

	testCases := map[string]struct {
		args                    args
		wantRefreshHappen       bool
		wantMasterNodesExtra    bool
		wantIsPeerInMasterNodes bool
	}{
		"no-refesh_IsPeerInMasterNodes-false": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC(),
				testInitMasterNodesExtra:       pastel.MasterNodes{masterNodes2},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{},
				testRefreshMasterNodesExtraErr: nil,
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       false,
			wantMasterNodesExtra:    false,
			wantIsPeerInMasterNodes: false,
		},
		"no-refesh_IsPeerInMasterNodes-true": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC(),
				testInitMasterNodesExtra:       pastel.MasterNodes{masterNodes1, masterNodes2},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{},
				testRefreshMasterNodesExtraErr: nil,
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       false,
			wantMasterNodesExtra:    false,
			wantIsPeerInMasterNodes: true,
		},
		"refresh-master-node-extra-failed_IsPeerInMasterNodes-false": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitMasterNodesExtra:       pastel.MasterNodes{},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{},
				testRefreshMasterNodesExtraErr: errors.New("test error"),
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       false,
			wantMasterNodesExtra:    true,
			wantIsPeerInMasterNodes: false,
		},
		"refresh-master-node-extra-failed_IsPeerInMasterNodes-true": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitMasterNodesExtra:       pastel.MasterNodes{masterNodes1},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{},
				testRefreshMasterNodesExtraErr: errors.New("test error"),
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       false,
			wantMasterNodesExtra:    true,
			wantIsPeerInMasterNodes: true,
		},
		"refresh-master-node-extra-success_IsPeerInMasterNodes-false": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitMasterNodesExtra:       pastel.MasterNodes{masterNodes1},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{masterNodes2},
				testRefreshMasterNodesExtraErr: nil,
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       true,
			wantMasterNodesExtra:    false,
			wantIsPeerInMasterNodes: false,
		},
		"refresh-master-node-extra-success_IsPeerInMasterNodes-true": {
			args: args{
				testInitLastMasterNodesRefresh: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitMasterNodesExtra:       pastel.MasterNodes{},
				testRefreshMasterNodesExtra:    pastel.MasterNodes{masterNodes1, masterNodes2},
				testRefreshMasterNodesExtraErr: nil,
				testPeerAuthInfo: &auth.PeerAuthInfo{
					PastelID: masterNodes1.ExtKey,
					Address:  masterNodes1.ExtAddress,
				},
			},
			wantRefreshHappen:       true,
			wantMasterNodesExtra:    true,
			wantIsPeerInMasterNodes: true,
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			// set up masternodes top & extra mock
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesExtra(tc.args.testRefreshMasterNodesExtra, tc.args.testRefreshMasterNodesExtraErr)

			helper := &AuthHelper{
				pastelClient:               pastelClientMock,
				secInfo:                    &alts.SecInfo{},
				masterNodes:                tc.args.testInitMasterNodesExtra,
				masterNodesRefreshDuration: testRefreshDuration,
				lastMasterNodesRefresh:     tc.args.testInitLastMasterNodesRefresh,
			}

			assert.Equal(t, tc.wantIsPeerInMasterNodes, helper.isPeerInMasterNodes(context.Background(), tc.args.testPeerAuthInfo))
			if tc.wantRefreshHappen {
				assert.NotEqual(t, tc.args.testInitLastMasterNodesRefresh.Unix(), helper.lastMasterNodesRefresh.Unix())
				assert.Equal(t, len(tc.args.testRefreshMasterNodesExtra), len(helper.masterNodes))
			} else {
				assert.Equal(t, tc.args.testInitLastMasterNodesRefresh.Unix(), helper.lastMasterNodesRefresh.Unix())
			}

			if tc.wantMasterNodesExtra {
				//mock assertion
				pastelClientMock.AssertExpectations(t)
				pastelClientMock.AssertMasterNodesExtra(1, mock.Anything)
			}

		})
	}
}

// TestAuthenticator_GenAuthInfo verify GenAuthInfo function
func TestAuthenticator_GenAuthInfo(t *testing.T) {
	type args struct {
		testSignErr       error
		testSignature     []byte
		testInitTimeStamp time.Time
		testInitSignature []byte
		testSecInfo       *alts.SecInfo
	}

	testRefreshDuration := 5 * time.Minute

	testCases := map[string]struct {
		args        args
		wantRefresh bool
		wantError   error
	}{
		"no-refresh_gen-success": {
			args: args{
				testSignErr:       nil,
				testSignature:     []byte("new-sig"),
				testInitTimeStamp: time.Now().UTC(),
				testInitSignature: []byte("old-sig"),
				testSecInfo:       &alts.SecInfo{},
			},
			wantRefresh: false,
			wantError:   nil,
		},
		"new-refresh_gen-success": {
			args: args{
				testSignErr:       nil,
				testSignature:     []byte("new-sig"),
				testInitTimeStamp: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitSignature: []byte("old-sig"),
				testSecInfo:       &alts.SecInfo{},
			},
			wantRefresh: true,
			wantError:   nil,
		},
		"new-refresh_gen-failed": {
			args: args{
				testSignErr:       errors.New("sign-error"),
				testSignature:     []byte("new-sig"),
				testInitTimeStamp: time.Now().UTC().Add(-testRefreshDuration - timestampMarginDuration),
				testInitSignature: []byte("old-sig"),
				testSecInfo:       &alts.SecInfo{},
			},
			wantRefresh: true,
			wantError:   errors.New("refresh auth"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			// set up masternodes top & extra mock
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.args.testSignature, tc.args.testSignErr)

			helper := &AuthHelper{
				pastelClient:    pastelClientMock,
				secInfo:         &alts.SecInfo{},
				timestamp:       tc.args.testInitTimeStamp,
				signature:       tc.args.testInitSignature,
				expiredDuration: testRefreshDuration,
				masterNodes:     pastel.MasterNodes{},
			}

			newAuthInfo, err := helper.GenAuthInfo(context.Background())

			if tc.wantError == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantError.Error()))
			}

			if tc.wantError == nil {
				if tc.wantRefresh {
					assert.NotEqual(t, tc.args.testInitTimeStamp.Unix(), helper.timestamp.Unix())
					assert.Equal(t, tc.args.testSignature, newAuthInfo.Signature)
				} else {
					assert.Equal(t, tc.args.testInitTimeStamp.Unix(), helper.timestamp.Unix())
					assert.Equal(t, tc.args.testInitSignature, newAuthInfo.Signature)
				}
			}
		})
	}
}
