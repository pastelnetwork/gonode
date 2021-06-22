package artworkregister

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task/state"
	stateTest "github.com/pastelnetwork/gonode/common/service/task/test"
	p2pTest "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/probe"
	probeTest "github.com/pastelnetwork/gonode/probe/test"
	nodeTest "github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskRun(t *testing.T) {
	t.Parallel()

	type args struct {
		taskID    string
		ctx       context.Context
		returnErr error
	}

	testCases := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "task run without error",
			args:      args{"1", context.Background(), nil},
			assertion: assert.NoError,
		}, {
			name:      "task run with error",
			args:      args{"2", context.Background(), fmt.Errorf("runAction error")},
			assertion: assert.Error,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnID(testCase.args.taskID).
				ListenOnSetStatusNotifyFunc().
				ListenOnCancel().
				ListenOnRunAction(testCase.args.returnErr)
			task := &Task{
				Task: clientTask.Task,
			}

			testCase.assertion(t, task.Run(testCase.args.ctx))
			clientTask.AssertExpectations(t)
			clientTask.AssertCancelCall(1)
			clientTask.AssertRunActionCall(1, mock.Anything)
		})
	}
}

func TestTaskSession(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx               context.Context
		requiredStatus    Status
		statusUpdate      Status
		isPrimary         bool
		requiredStatusErr error
	}

	testCases := []struct {
		name              string
		args              args
		assertion         assert.ErrorAssertionFunc
		newActionCall     int
		updateStatusCall  int
		assertExpectation bool
	}{
		{
			name: "primary node session mode",
			args: args{
				ctx:               context.Background(),
				requiredStatus:    StatusTaskStarted,
				statusUpdate:      StatusPrimaryMode,
				isPrimary:         true,
				requiredStatusErr: nil,
			},
			assertion:         assert.NoError,
			newActionCall:     1,
			updateStatusCall:  1,
			assertExpectation: true,
		}, {
			name: "secondary node session mode",
			args: args{
				ctx:               context.Background(),
				requiredStatus:    StatusTaskStarted,
				statusUpdate:      StatusSecondaryMode,
				isPrimary:         false,
				requiredStatusErr: nil,
			},
			assertion:         assert.NoError,
			newActionCall:     1,
			updateStatusCall:  1,
			assertExpectation: true,
		}, {
			name: "session error",
			args: args{
				ctx:               context.Background(),
				requiredStatus:    StatusTaskStarted,
				requiredStatusErr: fmt.Errorf("Task started"),
			},
			assertion:         assert.Error,
			newActionCall:     0,
			updateStatusCall:  0,
			assertExpectation: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx).
				ListenOnUpdateStatus()

			task := &Task{
				Task: clientTask,
			}

			testCase.assertion(t, task.Session(testCase.args.ctx, testCase.args.isPrimary))

			if testCase.assertExpectation {
				clientTask.AssertExpectations(t)
			}
			clientTask.AssertRequiredStatusCall(1, testCase.args.requiredStatus)
			clientTask.AssertNewActionCall(testCase.newActionCall, mock.AnythingOfType(stateTest.ActionFn))
			clientTask.AssertUpdateStatusCall(testCase.updateStatusCall, testCase.args.statusUpdate)
		})
	}
}

func TestTaskAcceptedNodes(t *testing.T) {
	t.Parallel()

	type contextFunc func() context.Context

	cancelContext := func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		return ctx
	}

	type args struct {
		requiredStatus    Status
		requiredStatusErr error
		subcribeStatus    *state.Status
		serverCtx         contextFunc
		ctx               contextFunc
	}

	testCases := []struct {
		name               string
		args               args
		assertion          assert.ErrorAssertionFunc
		assertExpectation  bool
		requiredStatusCall int
		newActionCall      int
		subcribeStatusCall int
	}{
		{
			name: "node is accepted",
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               context.Background,
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		}, {
			name: "serverCtx is timeout",
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         cancelContext,
				ctx:               context.Background,
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		}, {
			name: "ctx is timeout",
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               cancelContext,
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		},
		{
			name: "node is not primary mode",
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: fmt.Errorf("node is not primary"),
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               context.Background,
			},
			assertion:          assert.Error,
			assertExpectation:  false,
			requiredStatusCall: 1,
			newActionCall:      0,
			subcribeStatusCall: 0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx()).
				ListenOnSubscribeStatus(testCase.args.subcribeStatus)

			task := &Task{
				Task: clientTask.Task,
			}

			//only check error assertion. dont care about acceptedNodes result
			_, err := task.AcceptedNodes(testCase.args.serverCtx())

			testCase.assertion(t, err)

			if testCase.assertExpectation {
				clientTask.AssertExpectations(t)
			}

			clientTask.AssertRequiredStatusCall(testCase.requiredStatusCall, testCase.args.requiredStatus)
			clientTask.AssertNewActionCall(testCase.newActionCall, mock.AnythingOfType(stateTest.ActionFn))
			clientTask.AssertSubscribeStatusCall(testCase.subcribeStatusCall)
		})
	}
}

func TestTaskSessionNode(t *testing.T) {
	t.Parallel()

	type fields struct {
		acceptedNodes Nodes
	}

	type args struct {
		ctx               context.Context
		nodeID            string
		requiredStatus    Status
		requiredStatusErr error
		statusUpdate      Status
		nodes             pastel.MasterNodes
		retunNodeErr      error
	}

	testCases := []struct {
		name               string
		fields             fields
		args               args
		assertion          assert.ErrorAssertionFunc
		assertExpectation  bool
		requiredStatusCall int
		newActionCall      int
		updateStatusCall   int
		pasteNodesTopCall  int
	}{
		{
			name:   "session node 2 accepted and update task status as accepted",
			fields: fields{Nodes{&Node{ID: "2"}}},
			args: args{
				ctx:               context.Background(),
				nodeID:            "1",
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				statusUpdate:      StatusConnected,
				nodes: pastel.MasterNodes{
					pastel.MasterNode{ExtKey: "1", ExtAddress: "127.0.0.1"},
					pastel.MasterNode{ExtKey: "2", ExtAddress: "127.0.0.2"},
					pastel.MasterNode{ExtKey: "3", ExtAddress: "127.0.0.3"},
				},
				retunNodeErr: nil,
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			updateStatusCall:   1,
			pasteNodesTopCall:  1,
		}, {
			name:   "session node 1 accepted but connected nodes < defaultNumberConnectedNodes",
			fields: fields{Nodes{}},
			args: args{
				ctx:               context.Background(),
				nodeID:            "1",
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				statusUpdate:      StatusConnected,
				nodes: pastel.MasterNodes{
					pastel.MasterNode{ExtKey: "1", ExtAddress: "127.0.0.1"},
					pastel.MasterNode{ExtKey: "2", ExtAddress: "127.0.0.2"},
					pastel.MasterNode{ExtKey: "3", ExtAddress: "127.0.0.3"},
				},
				retunNodeErr: nil,
			},
			assertion:          assert.NoError,
			assertExpectation:  false,
			requiredStatusCall: 1,
			newActionCall:      1,
			updateStatusCall:   0,
			pasteNodesTopCall:  1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx).
				ListenOnUpdateStatus()

			pastelClient := test.NewMockClient(t)
			pastelClient.ListenOnMasterNodesTop(testCase.args.nodes, testCase.args.retunNodeErr)

			service := &Service{
				pastelClient: pastelClient,
				config:       NewConfig(),
			}
			task := &Task{
				Service:  service,
				Task:     clientTask.Task,
				accepted: testCase.fields.acceptedNodes,
			}

			testCase.assertion(t, task.SessionNode(testCase.args.ctx, testCase.args.nodeID))

			if testCase.assertExpectation {
				clientTask.AssertExpectations(t)
				pastelClient.AssertExpectations(t)
			}

			clientTask.AssertRequiredStatusCall(testCase.requiredStatusCall, testCase.args.requiredStatus)
			clientTask.AssertNewActionCall(testCase.newActionCall, mock.AnythingOfType(stateTest.ActionFn))
			clientTask.AssertUpdateStatusCall(testCase.updateStatusCall, testCase.args.statusUpdate)
			pastelClient.AssertMasterNodesTopCall(testCase.pasteNodesTopCall, mock.Anything)
		})
	}
}

func TestTaskConnectTo(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx               context.Context
		pastelID          string
		nodeID            string
		sessID            string
		requiredStatus    Status
		requiredStatusErr error
		statusUpdate      Status
		nodes             pastel.MasterNodes
		retunNodeErr      error
		returnConnectErr  error
		returnSessionErr  error
	}

	testCases := []struct {
		name               string
		args               args
		connectedTo        *Node
		assertion          assert.ErrorAssertionFunc
		assertExpectation  bool
		requiredStatusCall int
		newActionCall      int
		updateStatusCall   int
		pasteNodesTopCall  int
		sessionCall        int
	}{
		{
			name: "connect to node 1 as secondary node succeed",
			args: args{
				ctx:               context.Background(),
				pastelID:          "123",
				nodeID:            "1",
				sessID:            "sesid",
				requiredStatus:    StatusSecondaryMode,
				requiredStatusErr: nil,
				statusUpdate:      StatusConnected,
				nodes: pastel.MasterNodes{
					pastel.MasterNode{ExtKey: "1", ExtAddress: "127.0.0.1"},
					pastel.MasterNode{ExtKey: "2", ExtAddress: "127.0.0.2"},
					pastel.MasterNode{ExtKey: "3", ExtAddress: "127.0.0.3"},
				},
				retunNodeErr:     nil,
				returnConnectErr: nil,
				returnSessionErr: nil,
			},
			connectedTo:        &Node{ID: "1", Address: "127.0.0.1"},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			updateStatusCall:   1,
			pasteNodesTopCall:  1,
			sessionCall:        1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx).
				ListenOnUpdateStatus()

			pastelClient := test.NewMockClient(t)
			pastelClient.ListenOnMasterNodesTop(testCase.args.nodes, testCase.args.retunNodeErr)

			nodeClient := nodeTest.NewMockClient(t)
			nodeClient.
				ListenOnConnectCall(testCase.args.returnConnectErr).
				ListenOnRegisterArtworkCall().
				ListenOnSessionCall(testCase.args.returnSessionErr)

			service := &Service{
				pastelClient: pastelClient,
				nodeClient:   nodeClient.Client,
				config:       NewConfig(),
			}
			service.config.PastelID = testCase.args.pastelID

			task := &Task{
				Service: service,
				Task:    clientTask.Task,
			}

			testCase.assertion(t, task.ConnectTo(testCase.args.ctx, testCase.args.nodeID, testCase.args.sessID))

			if testCase.assertExpectation {
				assert.Equal(t, testCase.connectedTo.ID, task.connectedTo.ID)
				assert.Equal(t, testCase.connectedTo.Address, task.connectedTo.Address)

				clientTask.AssertExpectations(t)
				pastelClient.AssertExpectations(t)
				nodeClient.Client.AssertExpectations(t)
				nodeClient.Connection.AssertExpectations(t)
				nodeClient.RegisterArtwork.AssertExpectations(t)
			}

			clientTask.AssertRequiredStatusCall(testCase.requiredStatusCall, testCase.args.requiredStatus)
			clientTask.AssertNewActionCall(testCase.newActionCall, mock.AnythingOfType(stateTest.ActionFn))
			clientTask.AssertUpdateStatusCall(testCase.updateStatusCall, testCase.args.statusUpdate)
			pastelClient.AssertMasterNodesTopCall(testCase.pasteNodesTopCall, mock.Anything)
			nodeClient.AssertSessionCall(testCase.sessionCall, mock.Anything, testCase.args.pastelID, testCase.args.sessID)
		})
	}
}

func TestTaskProbeImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx               context.Context
		requiredStatus    Status
		requiredStatusErr error
		statusUpdate      Status
		fingerPrint       probe.Fingerprints
		fingerPrintErr    error
		storeID           string
		storeErr          error
	}

	testCases := []struct {
		name               string
		args               args
		assertion          assert.ErrorAssertionFunc
		requiredStatusCall int
		newActionCall      int
		updateStatusCall   int
		fingerPrintCall    int
		storeCall          int
	}{
		{
			name: "return valid fingerprint",
			args: args{
				ctx:               context.Background(),
				requiredStatus:    StatusConnected,
				requiredStatusErr: nil,
				statusUpdate:      StatusImageUploaded,
				fingerPrint:       probe.Fingerprints{probe.Fingerprint{1}},
				fingerPrintErr:    nil,
				storeID:           "1",
				storeErr:          nil,
			},
			assertion:          assert.NoError,
			requiredStatusCall: 1,
			newActionCall:      2,
			updateStatusCall:   1,
			storeCall:          1,
			fingerPrintCall:    1,
		},
	}

	t.Run("group", func(t *testing.T) {
		//create tmp file to store fake image file
		tmpFile, err := ioutil.TempFile("", "*.png")
		assert.NoError(t, err)

		err = tmpFile.Close()
		assert.NoError(t, err)

		defer os.Remove(tmpFile.Name())

		artworkFile, err := stateTest.NewTestImageFile(filepath.Dir(tmpFile.Name()), filepath.Base(tmpFile.Name()))
		assert.NoError(t, err)

		for _, testCase := range testCases {
			testCase := testCase

			t.Run(testCase.name, func(t *testing.T) {

				cancelCtx, cancel := context.WithCancel(testCase.args.ctx)
				cancel()
				clientTask := stateTest.NewMockTask(t)
				clientTask.
					ListenOnRequiredStatus(testCase.args.requiredStatusErr).
					ListenOnNewAction(cancelCtx).
					ListenOnUpdateStatus()

				tensorClient := probeTest.NewMockTensor(t)
				tensorClient.ListenOnFingerprints(testCase.args.fingerPrint, testCase.args.fingerPrintErr)

				p2pClient := p2pTest.NewMockClient(t)
				p2pClient.ListenOnStore(testCase.args.storeID, testCase.args.storeErr)

				service := &Service{
					probeTensor: tensorClient.Tensor,
					p2pClient:   p2pClient.Client,
					config:      NewConfig(),
				}

				task := &Task{
					Task:    clientTask.Task,
					Service: service,
				}

				got, err := task.ProbeImage(testCase.args.ctx, artworkFile)

				testCase.assertion(t, err)
				assert.NotNil(t, got)

				clientTask.AssertExpectations(t)
				clientTask.AssertRequiredStatusCall(testCase.requiredStatusCall, testCase.args.requiredStatus)
				clientTask.AssertNewActionCall(testCase.newActionCall, mock.AnythingOfType(stateTest.ActionFn))
				clientTask.AssertUpdateStatusCall(testCase.updateStatusCall, testCase.args.statusUpdate)

				tensorClient.AssertFingerprintsCall(testCase.fingerPrintCall, mock.Anything, mock.Anything)
				p2pClient.AssertStoreCall(testCase.storeCall, mock.Anything, mock.IsType([]byte{}))
			})
		}

	})
}

func TestTaskPastelNodeByExtKey(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		nodeID    string
		nodes     pastel.MasterNodes
		returnErr error
	}

	testCases := []struct {
		name      string
		args      args
		want      *Node
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "get valid masternode with id 2",
			args: args{
				ctx:    context.Background(),
				nodeID: "2",
				nodes: pastel.MasterNodes{
					pastel.MasterNode{ExtKey: "1", ExtAddress: "127.0.0.1"},
					pastel.MasterNode{ExtKey: "2", ExtAddress: "127.0.0.2"},
					pastel.MasterNode{ExtKey: "3", ExtAddress: "127.0.0.3"},
				},
				returnErr: nil,
			},
			want:      &Node{ID: "2", Address: "127.0.0.2"},
			assertion: assert.NoError,
		}, {
			name: "get not found masternode with id 4",
			args: args{
				ctx:    context.Background(),
				nodeID: "4",
				nodes: pastel.MasterNodes{
					pastel.MasterNode{ExtKey: "1", ExtAddress: "127.0.0.1"},
					pastel.MasterNode{ExtKey: "2", ExtAddress: "127.0.0.2"},
					pastel.MasterNode{ExtKey: "3", ExtAddress: "127.0.0.3"},
				},
				returnErr: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			client := test.NewMockClient(t)
			client.ListenOnMasterNodesTop(testCase.args.nodes, testCase.args.returnErr)

			service := &Service{
				pastelClient: client,
			}
			task := &Task{
				Service: service,
			}
			got, err := task.pastelNodeByExtKey(testCase.args.ctx, testCase.args.nodeID)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)
			client.AssertExpectations(t)
			client.AssertMasterNodesTopCall(1, mock.Anything)
		})
	}
}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *Service
	}

	testCases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create new task",
			args: args{},
			want: statusNames[StatusTaskStarted],
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, NewTask(testCase.args.service).Task.Status().String())
		})
	}
}
