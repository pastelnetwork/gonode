package artworkregister

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	stateMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
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

func newTestImageFile(fileName string) (*artwork.File, error) {
	err := test.CreateBlankImage("./"+fileName, 400, 400)
	if err != nil {
		return nil, err
	}

	imageStorage := artwork.NewStorage(fs.NewFileStorage("./"))
	imageFile := artwork.NewFile(imageStorage, fileName)
	return imageFile, nil
}

func TestTaskRun(t *testing.T) {
	t.Parallel()

	t.Skip()

	type fields struct {
		Ticket *Ticket
	}

	type args struct {
		taskID        string
		ctx           context.Context
		networkFee    float64
		masterNodes   pastel.MasterNodes
		primarySessID string
		pastelIDS     []string
		fingerPrint   []byte
		returnErr     error
	}

	testCases := []struct {
		fields          fields
		args            args
		assertion       assert.ErrorAssertionFunc
		numSessIDCall   int
		numUpdateStatus int
	}{
		{
			fields: fields{&Ticket{MaximumFee: 0.5}},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{Fee: 0.3, ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
					pastel.MasterNode{Fee: 0.4, ExtAddress: "127.0.0.1:4448", ExtKey: "4"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3", "4"},
				fingerPrint:   []byte("match"),
				returnErr:     nil,
			},
			assertion:       assert.NoError,
			numSessIDCall:   3,
			numUpdateStatus: 3,
		},
	}

	t.Run("group", func(t *testing.T) {
		//create tmp image file
		artworkFile, err := newTestImageFile("test.png")
		assert.NoError(t, err)

		defer func() {
			os.Remove("./test.png")
		}()

		for i, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				nodeClient := test.NewMockClient()
				nodeClient.
					ListenOnConnect(testCase.args.returnErr).
					ListenOnRegisterArtwork().
					ListenOnSession(testCase.args.returnErr).
					ListenOnConnectTo(testCase.args.returnErr).
					ListenOnSessID(testCase.args.primarySessID).
					ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
					ListenOnDone()

				// custom probe image listening call to get thumbnail file.
				// should remove the generated thumbnail file
				customProbeImage := func(ctx context.Context, image *artwork.File) []byte {
					err := image.Remove()
					assert.NoError(t, err)
					return testCase.args.fingerPrint
				}

				nodeClient.RegArtWorkMock.On("ProbeImage",
					mock.Anything,
					mock.IsType(&artwork.File{})).
					Return(customProbeImage, testCase.args.returnErr)

				pastelClientMock := pastelMock.NewMockClient()
				pastelClientMock.
					ListenOnStorageNetworkFee(testCase.args.networkFee, testCase.args.returnErr).
					ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr)

				service := &Service{
					pastelClient: pastelClientMock.ClientMock,
					nodeClient:   nodeClient.ClientMock,
					config:       NewConfig(),
				}

				taskClient := stateMock.NewMockTask()
				taskClient.
					ListenOnID(testCase.args.taskID).
					ListenOnUpdateStatus().
					ListenOnSetStatusNotifyFunc()

				ticket := testCase.fields.Ticket
				ticket.Image = artworkFile
				task := &Task{
					Task:    taskClient.TaskMock,
					Service: service,
					Ticket:  ticket,
				}

				//create context with timeout to automatically end process after 1 sec
				ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
				defer cancel()
				testCase.assertion(t, task.Run(ctx))

				taskClient.TaskMock.AssertExpectations(t)
				taskClient.TaskMock.AssertCalled(t, "ID")
				taskClient.TaskMock.AssertCalled(t, "UpdateStatus", mock.Anything)
				taskClient.TaskMock.AssertCalled(t, "SetStatusNotifyFunc", mock.Anything)
				taskClient.TaskMock.AssertNumberOfCalls(t, "UpdateStatus", testCase.numUpdateStatus)

				// //pastelClient mock assertion
				pastelClientMock.ClientMock.AssertExpectations(t)
				pastelClientMock.ClientMock.AssertCalled(t, "StorageNetworkFee", mock.Anything)
				pastelClientMock.ClientMock.AssertNumberOfCalls(t, "StorageNetworkFee", 1)

				// //nodeClient mock assertion
				nodeClient.ClientMock.AssertExpectations(t)
				nodeClient.ConnectionMock.AssertExpectations(t)
				nodeClient.RegArtWorkMock.AssertExpectations(t)
				nodeClient.RegArtWorkMock.AssertCalled(t, "AcceptedNodes", mock.Anything)
				nodeClient.RegArtWorkMock.AssertNumberOfCalls(t, "AcceptedNodes", 1)
				nodeClient.RegArtWorkMock.AssertCalled(t, "SessID")
				nodeClient.RegArtWorkMock.AssertNumberOfCalls(t, "SessID", testCase.numSessIDCall)
				nodeClient.RegArtWorkMock.AssertCalled(t, "Session", mock.Anything, false)
				nodeClient.RegArtWorkMock.AssertCalled(t, "ConnectTo", mock.Anything, mock.Anything, testCase.args.primarySessID)
				nodeClient.RegArtWorkMock.AssertCalled(t, "ProbeImage", mock.Anything, mock.IsType(&artwork.File{}))
			})

		}
	})

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
		acceptNodeErr   error
	}

	testCases := []struct {
		args          args
		want          []string
		assertion     assert.ErrorAssertionFunc
		numSessIDCall int
	}{
		{
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
				acceptNodeErr:   nil,
			},
			assertion:     assert.NoError,
			numSessIDCall: 6,
			want:          []string{"1:127.0.0.1", "2:127.0.0.2", "4:127.0.0.4", "7:127.0.0.7"},
		}, {
			args: args{
				ctx:          context.Background(),
				primaryIndex: 0,
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.1", "2"},
					{"127.0.0.1", "3"},
				},
				primaryPastelID: "1",
				primarySessID:   "xdxdf",
				pastelIDS:       []string{"2", "3"},
				returnErr:       nil,
				acceptNodeErr:   fmt.Errorf("primary node not accepted"),
			},
			assertion:     assert.Error,
			numSessIDCall: 2,
			want:          nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			//create new client mock
			nodeClient := test.NewMockClient()
			nodeClient.
				ListenOnConnect(testCase.args.returnErr).
				ListenOnRegisterArtwork().
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.acceptNodeErr)

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

func TestTaskIsSuitableStorageNetworkFee(t *testing.T) {
	t.Parallel()

	type fields struct {
		Ticket *Ticket
	}

	type args struct {
		ctx        context.Context
		networkFee float64
		returnErr  error
	}

	testCases := []struct {
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
				networkFee: 0.49,
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
				networkFee: 0.51,
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

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClient := pastelMock.NewMockClient()
			pastelClient.ListenOnStorageNetworkFee(testCase.args.networkFee, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClient.ClientMock,
			}

			task := &Task{
				Service: service,
				Ticket:  testCase.fields.Ticket,
			}

			got, err := task.isSuitableStorageFee(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//pastelClient mock assertion
			pastelClient.ClientMock.AssertExpectations(t)
			pastelClient.ClientMock.AssertCalled(t, "StorageNetworkFee", testCase.args.ctx)
			pastelClient.ClientMock.AssertNumberOfCalls(t, "StorageNetworkFee", 1)
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
		fields    fields
		args      args
		want      node.List
		assertion assert.ErrorAssertionFunc
	}{
		{
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
			args: args{
				ctx:       context.Background(),
				returnMn:  nil,
				returnErr: fmt.Errorf("connection timeout"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClient := pastelMock.NewMockClient()
			pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClient.ClientMock,
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
			pastelClient.ClientMock.AssertExpectations(t)
			pastelClient.ClientMock.AssertCalled(t, "MasterNodesTop", mock.Anything)
			pastelClient.ClientMock.AssertNumberOfCalls(t, "MasterNodesTop", 1)
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
		args args
		want *Task
	}{
		{
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
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			task := NewTask(testCase.args.service, testCase.args.Ticket)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Ticket, task.Ticket)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
