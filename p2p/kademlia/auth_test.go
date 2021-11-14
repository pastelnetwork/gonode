package kademlia

import (
	"context"
	"fmt"
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
	// func (ath *AuthHelper) tryRefreshMasternodeList(ctx context.Context) {
	// 	ath.masterNodesMtx.Lock()
	// 	defer ath.masterNodesMtx.Unlock()
	// 	if len(ath.masterNodes) == 0 || time.Now().After(ath.lastMasterNodesRefresh.Add(ath.masterNodesRefreshDuration)) {
	// 		if err := ath.unsafeRefreshMasternodeList(ctx); err != nil {
	// 			log.WithContext(ctx).WithError(err).Error("Update master node list failed")
	// 			return
	// 		}
	// 		ath.lastMasterNodesRefresh = time.Now()
	// 	}
	// }
	type args struct {
		helper                     *AuthHelper
		initLastMasterNodesRefresh time.Time
		initMasterNodesExtra       pastel.MasterNodes
		refreshMasterNodesExtra    pastel.MasterNodes
		refreshMasterNodesExtraErr error
		peerAuthInfo               *auth.PeerAuthInfo
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
				initLastMasterNodesRefresh: time.Now(),
				initMasterNodesExtra:       pastel.MasterNodes{masterNodes2},
				refreshMasterNodesExtra:    pastel.MasterNodes{},
				refreshMasterNodesExtraErr: nil,
				peerAuthInfo: &auth.PeerAuthInfo{
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
				initLastMasterNodesRefresh: time.Now(),
				initMasterNodesExtra:       pastel.MasterNodes{masterNodes1, masterNodes2},
				refreshMasterNodesExtra:    pastel.MasterNodes{},
				refreshMasterNodesExtraErr: nil,
				peerAuthInfo: &auth.PeerAuthInfo{
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
				initLastMasterNodesRefresh: time.Now().Add(-testRefreshDuration - timestampMarginDuration),
				initMasterNodesExtra:       pastel.MasterNodes{},
				refreshMasterNodesExtra:    pastel.MasterNodes{},
				refreshMasterNodesExtraErr: errors.New("test error"),
				peerAuthInfo: &auth.PeerAuthInfo{
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
				initLastMasterNodesRefresh: time.Now().Add(-testRefreshDuration - timestampMarginDuration),
				initMasterNodesExtra:       pastel.MasterNodes{masterNodes1},
				refreshMasterNodesExtra:    pastel.MasterNodes{},
				refreshMasterNodesExtraErr: errors.New("test error"),
				peerAuthInfo: &auth.PeerAuthInfo{
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
				initLastMasterNodesRefresh: time.Now().Add(-testRefreshDuration - timestampMarginDuration),
				initMasterNodesExtra:       pastel.MasterNodes{masterNodes1},
				refreshMasterNodesExtra:    pastel.MasterNodes{masterNodes2},
				refreshMasterNodesExtraErr: nil,
				peerAuthInfo: &auth.PeerAuthInfo{
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
				initLastMasterNodesRefresh: time.Now().Add(-testRefreshDuration - timestampMarginDuration),
				initMasterNodesExtra:       pastel.MasterNodes{},
				refreshMasterNodesExtra:    pastel.MasterNodes{masterNodes1, masterNodes2},
				refreshMasterNodesExtraErr: nil,
				peerAuthInfo: &auth.PeerAuthInfo{
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
			pastelClientMock.ListenOnMasterNodesExtra(tc.args.refreshMasterNodesExtra, tc.args.refreshMasterNodesExtraErr)

			tc.args.helper = &AuthHelper{
				pastelClient:               pastelClientMock,
				secInfo:                    &alts.SecInfo{},
				expiredDuration:            testRefreshDuration,
				masterNodes:                tc.args.initMasterNodesExtra,
				masterNodesRefreshDuration: testRefreshDuration,
				lastMasterNodesRefresh:     tc.args.initLastMasterNodesRefresh,
			}

			assert.Equal(t, tc.wantIsPeerInMasterNodes, tc.args.helper.isPeerInMasterNodes(context.Background(), tc.args.peerAuthInfo))
			if tc.wantRefreshHappen {
				assert.NotEqual(t, tc.args.initLastMasterNodesRefresh.Unix(), tc.args.helper.lastMasterNodesRefresh.Unix())
				assert.Equal(t, len(tc.args.refreshMasterNodesExtra), len(tc.args.helper.masterNodes))
			} else {
				assert.Equal(t, tc.args.initLastMasterNodesRefresh.Unix(), tc.args.helper.lastMasterNodesRefresh.Unix())
			}

			if tc.wantMasterNodesExtra {
				//mock assertion
				pastelClientMock.AssertExpectations(t)
				pastelClientMock.AssertMasterNodesExtra(1, mock.Anything)
			}

		})
	}
}
