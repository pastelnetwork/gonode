package artworkregister

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	stateTest "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/test"
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
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			args:      args{"1", context.Background(), nil},
			assertion: assert.NoError,
		}, {
			args:      args{"2", context.Background(), fmt.Errorf("runAction error")},
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
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
		ch                chan struct{}
		requiredStatus    Status
		statusUpdate      Status
		isPrimary         bool
		requiredStatusErr error
	}

	testCases := []struct {
		args              args
		assertion         assert.ErrorAssertionFunc
		newActionCall     int
		updateStatusCall  int
		assertExpectation bool
	}{
		{
			args: args{
				ctx:               context.Background(),
				ch:                make(chan struct{}),
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
			args: args{
				ctx:               context.Background(),
				ch:                make(chan struct{}),
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
			args: args{
				ctx:               context.Background(),
				ch:                make(chan struct{}),
				requiredStatus:    StatusTaskStarted,
				requiredStatusErr: fmt.Errorf("Task started"),
			},
			assertion:         assert.Error,
			newActionCall:     0,
			updateStatusCall:  0,
			assertExpectation: false,
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx, testCase.args.ch).
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
		ch                chan struct{}
	}

	testCases := []struct {
		args               args
		assertion          assert.ErrorAssertionFunc
		assertExpectation  bool
		requiredStatusCall int
		newActionCall      int
		subcribeStatusCall int
	}{
		{
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               context.Background,
				ch:                make(chan struct{}),
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		}, {
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         cancelContext,
				ctx:               context.Background,
				ch:                make(chan struct{}),
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		}, {
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: nil,
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               cancelContext,
				ch:                make(chan struct{}),
			},
			assertion:          assert.NoError,
			assertExpectation:  true,
			requiredStatusCall: 1,
			newActionCall:      1,
			subcribeStatusCall: 1,
		},
		{
			args: args{
				requiredStatus:    StatusPrimaryMode,
				requiredStatusErr: fmt.Errorf("node is not primary"),
				subcribeStatus:    state.NewStatus(StatusConnected),
				serverCtx:         context.Background,
				ctx:               context.Background,
				ch:                make(chan struct{}),
			},
			assertion:          assert.Error,
			assertExpectation:  false,
			requiredStatusCall: 1,
			newActionCall:      0,
			subcribeStatusCall: 0,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx(), testCase.args.ch).
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

	type args struct {
		ctx               context.Context
		nodeID            string
		requiredStatus    Status
		requiredStatusErr error
		ch                chan struct{}
		statusUpdate      Status
		nodes             pastel.MasterNodes
		retunNodeErr      error
	}

	testCases := []struct {
		args               args
		assertion          assert.ErrorAssertionFunc
		assertExpectation  bool
		requiredStatusCall int
		newActionCall      int
		updateStatusCall   int
		pasteNodesTopCall  int
	}{
		// TODO: Add test cases.
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			clientTask := stateTest.NewMockTask(t)
			clientTask.
				ListenOnRequiredStatus(testCase.args.requiredStatusErr).
				ListenOnNewAction(testCase.args.ctx, testCase.args.ch).
				ListenOnUpdateStatus()

			pastelClient := test.NewMockClient(t)
			pastelClient.ListenOnMasterNodesTop(testCase.args.nodes, testCase.args.retunNodeErr)

			service := &Service{
				pastelClient: pastelClient,
			}
			task := &Task{
				Service: service,
				Task:    clientTask.Task,
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
	type fields struct {
		Task             task.Task
		Service          *Service
		ResampledArtwork *artwork.File
		Artwork          *artwork.File
		acceptedMu       sync.Mutex
		accpeted         Nodes
		connectedTo      *Node
	}
	type args struct {
		in0    context.Context
		nodeID string
		sessID string
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
				Task:             tt.fields.Task,
				Service:          tt.fields.Service,
				ResampledArtwork: tt.fields.ResampledArtwork,
				Artwork:          tt.fields.Artwork,
				acceptedMu:       tt.fields.acceptedMu,
				accpeted:         tt.fields.accpeted,
				connectedTo:      tt.fields.connectedTo,
			}
			tt.assertion(t, task.ConnectTo(tt.args.in0, tt.args.nodeID, tt.args.sessID))
		})
	}
}

func TestTaskProbeImage(t *testing.T) {
	type fields struct {
		Task             task.Task
		Service          *Service
		ResampledArtwork *artwork.File
		Artwork          *artwork.File
		acceptedMu       sync.Mutex
		accpeted         Nodes
		connectedTo      *Node
	}
	type args struct {
		in0  context.Context
		file *artwork.File
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []byte
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:             tt.fields.Task,
				Service:          tt.fields.Service,
				ResampledArtwork: tt.fields.ResampledArtwork,
				Artwork:          tt.fields.Artwork,
				acceptedMu:       tt.fields.acceptedMu,
				accpeted:         tt.fields.accpeted,
				connectedTo:      tt.fields.connectedTo,
			}
			got, err := task.ProbeImage(tt.args.in0, tt.args.file)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
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
		args      args
		want      *Node
		assertion assert.ErrorAssertionFunc
	}{
		{
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

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
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
		args args
		want string
	}{
		{
			args: args{},
			want: statusNames[StatusTaskStarted],
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.want, NewTask(testCase.args.service).Task.Status().String())
		})
	}
}
