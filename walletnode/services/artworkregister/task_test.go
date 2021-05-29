package artworkregister

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNode(address, pastelID string) *Node {
	return &Node{Address: address, PastelID: pastelID}
}

func pullPastelIDNodes(nodes Nodes) []string {
	var v []string
	for _, n := range nodes {
		v = append(v, n.PastelID)
	}

	sort.Strings(v)
	return v
}

func TestTaskRun(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			tt.assertion(t, task.Run(tt.args.ctx))
		})
	}
}

func TestTask_run(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			tt.assertion(t, task.run(tt.args.ctx))
		})
	}
}

func TestTaskMeshNodes(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx             context.Context
		nodes           Nodes
		primaryIndex    int
		primaryPastelID string
		primarySessID   string
		pastelIDS       []string
		returnErr       error
	}
	tests := []struct {
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
				nodes: Nodes{
					newTestNode("127.0.0.1", "1"),
					newTestNode("127.0.0.2", "2"),
					newTestNode("127.0.0.3", "3"),
					newTestNode("127.0.0.4", "4"),
					newTestNode("127.0.0.5", "5"),
					newTestNode("127.0.0.6", "6"),
					newTestNode("127.0.0.7", "7"),
				},
				primaryPastelID: "2",
				primarySessID:   "xdcfjc",
				pastelIDS:       []string{"1", "4", "7"},
				returnErr:       nil,
			},
			assertion:     assert.NoError,
			numSessIDCall: 6,
			want:          []string{"1", "2", "4", "7"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			//setup mock service for each node
			nodesClientMock := []*test.Client{}
			nodes := tt.args.nodes

			for _, node := range nodes {
				nodeClient := test.NewMockClient()
				nodeClient.
					ListenOnConnect(tt.args.returnErr).
					ListenOnRegisterArtwork().
					ListenOnSession(tt.args.returnErr).
					ListenOnConnectTo(tt.args.returnErr).
					ListenOnSessID(tt.args.primarySessID).
					//blocking primary acceptNode call with delay 2 second
					ListenOnAcceptedNodes(2, tt.args.pastelIDS, tt.args.returnErr)

				node.client = nodeClient.ClientMock
				nodesClientMock = append(nodesClientMock, nodeClient)
			}

			task := &Task{}
			got, err := task.meshNodes(tt.args.ctx, tt.args.nodes, tt.args.primaryIndex)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, pullPastelIDNodes(got))

			//primary node mock should call these methods
			primaryNodeMock := nodesClientMock[tt.args.primaryIndex]
			primaryNodeMock.RegArtWorkMock.AssertCalled(t, "AcceptedNodes", mock.Anything)
			primaryNodeMock.RegArtWorkMock.AssertNumberOfCalls(t, "AcceptedNodes", 1)
			primaryNodeMock.RegArtWorkMock.AssertCalled(t, "SessID")
			primaryNodeMock.RegArtWorkMock.AssertNumberOfCalls(t, "SessID", tt.numSessIDCall)

			t.Run("group", func(t *testing.T) {
				for i, n := range nodesClientMock {
					i, n := i, n

					t.Run(fmt.Sprintf("testCall-%d", i), func(t *testing.T) {
						n.ClientMock.AssertExpectations(t)
						n.ConnectionMock.AssertExpectations(t)

						//if node not primary should call these method
						if i != tt.args.primaryIndex {
							n.RegArtWorkMock.AssertCalled(t, "Session", mock.Anything, false)
							n.RegArtWorkMock.AssertCalled(t, "ConnectTo", mock.Anything, tt.args.primaryPastelID, tt.args.primarySessID)
						}
					})
				}
			})

		})
	}

}

func TestTaskIsSuitableStorageFee(t *testing.T) {
	t.Parallel()

	type fields struct {
		Ticket *Ticket
	}

	type args struct {
		ctx       context.Context
		pastelFee *pastel.StorageFee
		returnErr error
	}

	tests := []struct {
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
				ctx:       context.Background(),
				pastelFee: &pastel.StorageFee{NetworkFee: 0.49},
			},
			want:      true,
			assertion: assert.NoError,
		},
		{
			fields: fields{
				Ticket: &Ticket{MaximumFee: 0.5},
			},
			args: args{
				ctx:       context.Background(),
				pastelFee: &pastel.StorageFee{NetworkFee: 0.51},
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

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClientMock := test.NewMockPastelClient()
			pastelClientMock.ListenOnStorageFee(tt.args.pastelFee, tt.args.returnErr)
			service := &Service{
				pastelClient: pastelClientMock.ClientMock,
			}

			task := &Task{
				Service: service,
				Ticket:  tt.fields.Ticket,
			}

			got, err := task.isSuitableStorageFee(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			//pastelClient mock assertion
			pastelClientMock.ClientMock.AssertExpectations(t)
			pastelClientMock.ClientMock.AssertCalled(t, "StorageFee", tt.args.ctx)
			pastelClientMock.ClientMock.AssertNumberOfCalls(t, "StorageFee", 1)
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

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      Nodes
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
			want: Nodes{
				&Node{Address: "127.0.0.1:4444", PastelID: "1"},
				&Node{Address: "127.0.0.1:4445", PastelID: "2"},
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
			want: Nodes{
				&Node{Address: "127.0.0.1:4445", PastelID: "2"},
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

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClientMock := test.NewMockPastelClient()
			pastelClientMock.ListenOnMasterNodesTop(tt.args.returnMn, tt.args.returnErr)
			service := &Service{
				pastelClient: pastelClientMock.ClientMock,
			}

			task := &Task{
				Task:    tt.fields.Task,
				Service: service,
				Ticket:  tt.fields.Ticket,
			}
			got, err := task.pastelTopNodes(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			//mock assertion
			pastelClientMock.ClientMock.AssertExpectations(t)
			pastelClientMock.ClientMock.AssertCalled(t, "MasterNodesTop", mock.Anything)
			pastelClientMock.ClientMock.AssertNumberOfCalls(t, "MasterNodesTop", 1)
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

	tests := []struct {
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
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(tt.args.service, tt.args.Ticket)
			assert.Equal(t, tt.want.Service, task.Service)
			assert.Equal(t, tt.want.Ticket, task.Ticket)
			assert.Equal(t, tt.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
