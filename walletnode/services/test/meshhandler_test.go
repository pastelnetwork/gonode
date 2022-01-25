package test

import (
	"fmt"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/sense_register"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	service "github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"testing"
)

func TestNodesDisconnectInactive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		conn []struct {
			client       *test.Client
			activated    bool
			numberOfCall int
		}
	}{
		{
			conn: []struct {
				client       *test.Client
				activated    bool
				numberOfCall int
			}{
				{
					client:       test.NewMockClient(t),
					numberOfCall: 1,
					activated:    false,
				},
				{
					client:       test.NewMockClient(t),
					numberOfCall: 0,
					activated:    true,
				},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			meshHandler := mixins.NewMeshHandlerSimple(test.NewMockClient(t), service.RegisterSenseNodeMaker{})

			nodes := common.SuperNodeList{}

			for _, c := range testCase.conn {
				c.client.ListenOnClose(nil)

				node := common.NewSuperNode(nil, "", "", nil)
				node.ConnectionInterface = c.client.Connection
				node.SetActive(c.activated)
				nodes = append(nodes, node)
			}

			meshHandler.Nodes = nodes

			meshHandler.DisconnectInactiveNodes(ctx)

			for j, c := range testCase.conn {
				c := c

				t.Run(fmt.Sprintf("close-called-%d", j), func(t *testing.T) {
					c.client.AssertCloseCall(c.numberOfCall)
				})

			}

		})
	}

}
