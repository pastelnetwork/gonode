package artworkregister

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	pastelMock "github.com/pastelnetwork/gonode/walletnode/node/test/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNode(address, pastelID string) *node.Node {
	return node.NewNode(nil, address, pastelID)
}

func pullPastelAddressIDNodes(nodes node.List) []string {
	var v []string
	for _, n := range nodes {
		v = append(v, fmt.Sprintf("%s:%s", n.PastelID(), n.String()))
	}

	sort.Strings(v)
	return v
}

func TestTaskMeshNodes(t *testing.T) {
	t.Parallel()

	type nodeArg struct {
		address  string
		pastelID string
	}

	type args struct {
		ctx             context.Context
		nodes           []nodeArg
		primaryIndex    int
		primaryPastelID string
		primarySessID   string
		pastelIDS       []string
		returnErr       error
	}

	testCases := []struct {
		name          string
		args          args
		want          []string
		assertion     assert.ErrorAssertionFunc
		numSessIDCall int
	}{
		{
			name: "return 3 secondary nodes",
			args: args{
				ctx:          context.Background(),
				primaryIndex: 1,
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
					{"127.0.0.3", "3"},
					{"127.0.0.4", "4"},
					{"127.0.0.5", "5"},
					{"127.0.0.6", "6"},
					{"127.0.0.7", "7"},
				},
				primaryPastelID: "2",
				primarySessID:   "xdcfjc",
				pastelIDS:       []string{"1", "4", "7"},
				returnErr:       nil,
			},
			assertion:     assert.NoError,
			numSessIDCall: 6,
			want:          []string{"1:127.0.0.1", "2:127.0.0.2", "4:127.0.0.4", "7:127.0.0.7"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//create new client mock
			nodeClient := test.NewMockClient()
			nodeClient.
				ListenOnConnect(testCase.args.returnErr).
				ListenOnRegisterArtwork().
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr)

			nodes := node.List{}
			for _, n := range testCase.args.nodes {
				nodes.Add(node.NewNode(nodeClient.ClientMock, n.address, n.pastelID))
			}

			task := &Task{}
			got, err := task.meshNodes(testCase.args.ctx, nodes, testCase.args.primaryIndex)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, pullPastelAddressIDNodes(got))

			nodeClient.RegArtWorkMock.AssertCalled(t, "AcceptedNodes", mock.Anything)
			nodeClient.RegArtWorkMock.AssertNumberOfCalls(t, "AcceptedNodes", 1)
			nodeClient.RegArtWorkMock.AssertCalled(t, "SessID")
			nodeClient.RegArtWorkMock.AssertNumberOfCalls(t, "SessID", testCase.numSessIDCall)
			nodeClient.RegArtWorkMock.AssertCalled(t, "Session", mock.Anything, false)
			nodeClient.RegArtWorkMock.AssertCalled(t, "ConnectTo", mock.Anything, testCase.args.primaryPastelID, testCase.args.primarySessID)

			nodeClient.ClientMock.AssertExpectations(t)
			nodeClient.ConnectionMock.AssertExpectations(t)

		})

	}
}

func TestTaskIsSuitableStorageFee(t *testing.T) {
	t.Parallel()

	type fields struct {
		Ticket *Ticket
	}

	type args struct {
		ctx        context.Context
		networkFee *pastel.StorageFee
		returnErr  error
	}

	testCases := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				Ticket: &Ticket{MaximumFee: 0.5},
			},
			args: args{
				ctx:        context.Background(),
				networkFee: &pastel.StorageFee{NetworkFee: 0.49},
			},
			want:      true,
			assertion: assert.NoError,
		},
		{
			fields: fields{
				Ticket: &Ticket{MaximumFee: 0.5},
			},
			args: args{
				ctx:        context.Background(),
				networkFee: &pastel.StorageFee{NetworkFee: 0.51},
			},
			want:      false,
			assertion: assert.NoError,
		},
		{
			args: args{
				ctx:       context.Background(),
				returnErr: fmt.Errorf("connection timeout"),
			},
			want:      false,
			assertion: assert.Error,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClientMock := pastelMock.NewMockPastelClient()
			pastelClientMock.ListenOnStorageFee(testCase.args.networkFee, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClientMock.PastelMock,
			}

			task := &Task{
				Service: service,
				Ticket:  testCase.fields.Ticket,
			}

			got, err := task.isSuitableStorageFee(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//pastelClient mock assertion
			pastelClientMock.PastelMock.AssertExpectations(t)
			pastelClientMock.PastelMock.AssertCalled(t, "StorageFee", testCase.args.ctx)
			pastelClientMock.PastelMock.AssertNumberOfCalls(t, "StorageFee", 1)
		})
	}
}

func TestTaskPastelTopNodes(t *testing.T) {
	t.Parallel()

	type fields struct {
		Task   task.Task
		Ticket *Ticket
	}

	type args struct {
		ctx       context.Context
		returnMn  pastel.MasterNodes
		returnErr error
	}

	testCases := []struct {
		name      string
		fields    fields
		args      args
		want      node.List
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "valid fee",
			fields: fields{
				Ticket: &Ticket{
					MaximumFee: 0.30,
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4444", "1"),
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			name: "only one node valid fee",
			fields: fields{
				Ticket: &Ticket{
					MaximumFee: 0.3,
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.5, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			name: "failed retrieve top master node list",
			args: args{
				ctx:       context.Background(),
				returnMn:  nil,
				returnErr: fmt.Errorf("connection timeout"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClientMock := pastelMock.NewMockPastelClient()
			pastelClientMock.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClientMock.PastelMock,
			}

			task := &Task{
				Task:    testCase.fields.Task,
				Service: service,
				Ticket:  testCase.fields.Ticket,
			}
			got, err := task.pastelTopNodes(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//mock assertion
			pastelClientMock.PastelMock.AssertExpectations(t)
			pastelClientMock.PastelMock.AssertCalled(t, "MasterNodesTop", mock.Anything)
			pastelClientMock.PastelMock.AssertNumberOfCalls(t, "MasterNodesTop", 1)
		})
	}

}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *Service
		Ticket  *Ticket
	}

	service := &Service{}
	ticket := &Ticket{}

	testCases := []struct {
		name string
		args args
		want *Task
	}{
		{
			name: "create new task",
			args: args{
				service: service,
				Ticket:  ticket,
			},
			want: &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
				Ticket:  ticket,
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			task := NewTask(testCase.args.service, testCase.args.Ticket)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Ticket, task.Ticket)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
