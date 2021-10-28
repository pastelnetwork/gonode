package kademlia

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestConfigureBootstrapNodes(t *testing.T) {
	type args struct {
		s                   *DHT
		badAddrs            []string
		masterNodes         pastel.MasterNodes
		masterNodesExtra    pastel.MasterNodes
		masterNodesErr      error
		masterNodesExtraErr error
	}

	testCases := map[string]struct {
		args                  args
		wantErr               error
		wantBootstrapNodesLen int
		wantNodeIP            string
		wantNodePort          int
	}{
		"success": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 0,
		},
		"success-top": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "10.0.10.1:1234"},
				},
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 1,
			wantNodeIP:            "10.0.10.1",
			wantNodePort:          1234,
		},
		"success-extra": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
				masterNodes:         pastel.MasterNodes{},
				masterNodesExtra: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "10.0.10.1:1111"},
				},
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 1,
			wantNodeIP:            "10.0.10.1",
			wantNodePort:          1111,
		},
		"success-ignore-bad-addr": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
				masterNodes:         pastel.MasterNodes{},
				masterNodesExtra: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1111"},
					pastel.MasterNode{ExtP2P: "1.1.1.1:1313"},
				},
				badAddrs: []string{"0.0.0.0:1111"},
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 1,
			wantNodeIP:            "1.1.1.1",
			wantNodePort:          1313,
		},
		"success-ignore-bad-addr-top": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1234"},
					pastel.MasterNode{ExtP2P: "0.0.0.0:8784"},
				},
				masterNodesExtra: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1111"},
					pastel.MasterNode{ExtP2P: "1.1.1.1:1919"},
				},
				badAddrs: []string{"0.0.0.0:1234", "0.0.0.0:1111", "0.0.0.0:8784"},
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 1,
			wantNodeIP:            "1.1.1.1",
			wantNodePort:          1919,
		},
		"success-ignore-bad-addr-all": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: nil,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1234"},
					pastel.MasterNode{ExtP2P: "0.0.0.0:8784"},
				},
				masterNodesExtra: pastel.MasterNodes{
					pastel.MasterNode{ExtP2P: "0.0.0.0:1111"},
					pastel.MasterNode{ExtP2P: "1.1.1.1:1919"},
				},
				badAddrs: []string{"0.0.0.0:1234", "0.0.0.0:1111", "0.0.0.0:8784", "1.1.1.1:1919"},
			},
			wantErr:               nil,
			wantBootstrapNodesLen: 0,
		},
		"top-err": {
			args: args{
				masterNodesErr:      errors.New("test"),
				masterNodesExtraErr: nil,
			},
			wantErr:               errors.New("test"),
			wantBootstrapNodesLen: 0,
		},
		"extra-err": {
			args: args{
				masterNodesErr:      nil,
				masterNodesExtraErr: errors.New("test"),
			},
			wantErr:               errors.New("test"),
			wantBootstrapNodesLen: 0,
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			// set up masternodes top & extra mock
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(tc.args.masterNodes, tc.args.masterNodesErr).
				ListenOnMasterNodesExtra(tc.args.masterNodesExtra, tc.args.masterNodesExtraErr)

			tc.args.s = &DHT{
				options:      &Options{BootstrapNodes: []*Node{}},
				pastelClient: pastelClientMock,
				done:         make(chan struct{}),
				cache:        memory.NewKeyValue(),
			}

			for _, addr := range tc.args.badAddrs {
				err := tc.args.s.cache.Set(addr, []byte("true"))
				assert.Nil(t, err)
			}

			err := tc.args.s.ConfigureBootstrapNodes(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantBootstrapNodesLen, len(tc.args.s.options.BootstrapNodes))
				if tc.wantBootstrapNodesLen == 1 {
					assert.Equal(t, tc.wantNodeIP, tc.args.s.options.BootstrapNodes[0].IP)
					assert.Equal(t, tc.wantNodePort, tc.args.s.options.BootstrapNodes[0].Port)
				}
			}
		})
	}
}
