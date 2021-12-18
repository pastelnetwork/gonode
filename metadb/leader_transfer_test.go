package metadb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
	"github.com/pastelnetwork/gonode/pastel"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
)

func TestInitLeadershipTransferTrigger(t *testing.T) {
	type args struct {
		leaderCheckInterval     time.Duration
		blockCountCheckInterval time.Duration
		retCount                int32
		masterNodes             pastel.MasterNodes
		masterNodesExtra        pastel.MasterNodes
		masterNodesExtraErr     error
		masterNodesErr          error
		blockCountErr           error
	}

	testCases := map[string]struct {
		service *service
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				leaderCheckInterval:     time.Microsecond,
				blockCountCheckInterval: time.Microsecond,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1234", IPAddress: "127.0.0.1:9090"},
					pastel.MasterNode{ExtP2P: "0.0.0.0:8784", IPAddress: "127.0.0.1:9190"},
				},
				masterNodesExtra: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1111", IPAddress: "127.0.0.1:9090"},
					pastel.MasterNode{ExtP2P: "1.1.1.1:1919", IPAddress: "127.0.0.1:9190"},
				},
				masterNodesErr: nil,
				blockCountErr:  nil,
			},
			service: &service{
				config: NewConfig(),
				nodeID: "node-id",
			},
			wantErr: nil,
		},
		"mn-error": {
			args: args{
				leaderCheckInterval:     time.Microsecond,
				blockCountCheckInterval: time.Microsecond,
				masterNodesErr:          errors.New("test"),
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1234", IPAddress: "127.0.0.1:9090"},
					pastel.MasterNode{ExtP2P: "0.0.0.0:8784", IPAddress: "127.0.0.1:9190"},
				},
			},
			service: &service{
				config: NewConfig(),
				nodeID: "node-id",
			},
		},
		"block_count-error": {
			args: args{
				leaderCheckInterval:     time.Microsecond,
				blockCountCheckInterval: time.Microsecond,
				masterNodesErr:          nil,
				blockCountErr:           errors.New("test"),
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1234", IPAddress: "127.0.0.1:9090"},
					pastel.MasterNode{ExtP2P: "0.0.0.0:8784", IPAddress: "127.0.0.1:9190"},
				},
			},
			service: &service{
				config: NewConfig(),
				nodeID: "node-id",
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(tc.args.retCount, tc.args.blockCountErr)
			pastelClientMock.ListenOnMasterNodesTop(tc.args.masterNodes, tc.args.masterNodesErr).
				ListenOnMasterNodesExtra(tc.args.masterNodesExtra, tc.args.masterNodesExtraErr)

			tc.service.pastelClient = pastelClientMock

			s := store.MustNewStoreTest(true)
			defer os.RemoveAll(s.Path())

			if err := s.Open(true); err != nil {
				t.Fatalf("failed to open single-node store: %s", err.Error())
			}

			_, err := s.WaitForLeader(context.TODO(), 10*time.Second)
			if err != nil {
				t.Fatalf("Error waiting for leader: %s", err)
			}
			_, err = s.LeaderAddr()
			if err != nil {
				t.Fatalf("failed to get leader address: %s", err.Error())
			}
			tc.service.db = s

			tc.service.initLeadershipTransferTrigger(ctx, tc.args.leaderCheckInterval,
				tc.args.blockCountCheckInterval)
			time.Sleep(2 * time.Second)
			cancelFunc()
		})
	}
}

func TestTransferLeadership(t *testing.T) {
	type args struct{}

	testCases := map[string]struct {
		service *service
		args    args
		wantErr error
	}{
		"success": {
			args: args{},
			service: &service{
				config: NewConfig(),
				nodeID: "node-id",
			},
			wantErr: nil,
		},
		"err-notleader": {
			args: args{},
			service: &service{
				config: NewConfig(),
				nodeID: "node-id",
			},
			wantErr: errors.New(store.ErrNotLeader),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			s0 := store.MustNewStoreTest(true)
			defer os.RemoveAll(s0.Path())

			if err := s0.Open(true); err != nil {
				t.Fatalf("failed to open node for multi-node test: %s", err.Error())
			}
			defer s0.Close(true)

			if _, err := s0.WaitForLeader(context.TODO(), 10*time.Second); err != nil {
				t.Fatalf("Error waiting for leader: %s", err)
			}

			s1 := store.MustNewStoreTest(true)
			defer os.RemoveAll(s1.Path())
			if err := s1.Open(false); err != nil {
				t.Fatalf("failed to open node for multi-node test: %s", err.Error())
			}
			defer s1.Close(true)

			// Join the second node to the first.
			if err := s0.Join(s1.ID(), s1.Addr(), true); err != nil {
				t.Fatalf("failed to join to node at %s: %s", s0.Addr(), err.Error())
			}

			got, err := s1.WaitForLeader(context.TODO(), 10*time.Second)
			if err != nil {
				t.Fatalf("failed to get leader address on follower: %s", err.Error())
			}

			// Check leader state on follower.
			if exp := s0.Addr(); got != exp {
				t.Fatalf("wrong leader address returned, got: %s, exp %s", got, exp)
			}

			if tc.wantErr == nil {
				tc.service.db = s0
				err = tc.service.transferLeadership(context.Background(), 10, s1.ID(), s1.Addr())
			} else {
				tc.service.db = s1
				err = tc.service.transferLeadership(context.Background(), 10, s0.ID(), s0.Addr())
			}

			want := "-"
			if tc.wantErr != nil {
				want = tc.wantErr.Error()
			}

			have := "-"
			if err != nil {
				have = err.Error()
			}
			assert.Equal(t, want, have)

			if tc.wantErr == nil {
				newLeader, err := s0.WaitForLeader(context.Background(), 2*time.Second)
				assert.Nil(t, err)

				assert.Equal(t, s1.Addr(), newLeader)
			}
		})
	}
}
