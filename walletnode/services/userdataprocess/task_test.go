package userdataprocess

//
//import (
//	"context"
//	"fmt"
//	"reflect"
//	"sort"
//	"testing"
//
//	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
//	"github.com/pastelnetwork/gonode/common/service/task"
//	"github.com/pastelnetwork/gonode/common/service/userdata"
//	"github.com/pastelnetwork/gonode/common/types"
//	"github.com/pastelnetwork/gonode/pastel"
//	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
//	nodeTest "github.com/pastelnetwork/gonode/walletnode/node/test/process_userdata"
//	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess/node"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//func newTestNode(address, pastelID string) *node.Node {
//	return node.NewNode(nil, address, pastelID)
//}
//
//func pullPastelAddressIDNodes(nodes node.List) []string {
//	var v []string
//	for _, n := range nodes {
//		v = append(v, fmt.Sprintf("%s:%s", n.PastelID(), n.String()))
//	}
//
//	sort.Strings(v)
//	return v
//}
//
//func TestTask_Run(t *testing.T) {
//	type fields struct {
//		Task          task.Task
//		Service       *UserDataService
//		nodes         node.List
//		resultChan    chan *userdata.ProcessResult
//		err           error
//		request       *userdata.ProcessRequest
//		userpastelid  string
//		resultChanGet chan *userdata.ProcessRequest
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			task := &UserDataTask{
//				Task:            tt.fields.Task,
//				UserDataService: tt.fields.Service,
//				nodes:           tt.fields.nodes,
//				resultChan:      tt.fields.resultChan,
//				err:             tt.fields.err,
//				request:         tt.fields.request,
//				userpastelid:    tt.fields.userpastelid,
//				resultChanGet:   tt.fields.resultChanGet,
//			}
//			if err := task.Run(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("Task.Run() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestTask_run(t *testing.T) {
//	type fields struct {
//		Task          task.Task
//		Service       *UserDataService
//		nodes         node.List
//		resultChan    chan *userdata.ProcessResult
//		err           error
//		request       *userdata.ProcessRequest
//		userpastelid  string
//		resultChanGet chan *userdata.ProcessRequest
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			task := &UserDataTask{
//				Task:            tt.fields.Task,
//				UserDataService: tt.fields.Service,
//				nodes:           tt.fields.nodes,
//				resultChan:      tt.fields.resultChan,
//				err:             tt.fields.err,
//				request:         tt.fields.request,
//				userpastelid:    tt.fields.userpastelid,
//				resultChanGet:   tt.fields.resultChanGet,
//			}
//			if err := task.run(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("Task.run() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestTask_AggregateResult(t *testing.T) {
//	type fields struct {
//		Task          task.Task
//		Service       *UserDataService
//		nodes         node.List
//		resultChan    chan *userdata.ProcessResult
//		err           error
//		request       *userdata.ProcessRequest
//		userpastelid  string
//		resultChanGet chan *userdata.ProcessRequest
//	}
//	type args struct {
//		ctx   context.Context
//		nodes node.List
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    userdata.ProcessResult
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			task := &UserDataTask{
//				Task:            tt.fields.Task,
//				UserDataService: tt.fields.Service,
//				nodes:           tt.fields.nodes,
//				resultChan:      tt.fields.resultChan,
//				err:             tt.fields.err,
//				request:         tt.fields.request,
//				userpastelid:    tt.fields.userpastelid,
//				resultChanGet:   tt.fields.resultChanGet,
//			}
//			got, err := task.AggregateResult(tt.args.ctx, tt.args.nodes)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Task.AggregateResult() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Task.AggregateResult() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestTask_meshNodes(t *testing.T) {
//	t.Parallel()
//
//	type nodeArg struct {
//		address  string
//		pastelID string
//	}
//	type args struct {
//		ctx             context.Context
//		nodes           []nodeArg
//		primaryIndex    int
//		primaryPastelID string
//		primarySessID   string
//		pastelIDS       []string
//		returnErr       error
//		acceptNodeErr   error
//	}
//	tests := []struct {
//		args             args
//		want             []string
//		assertion        assert.ErrorAssertionFunc
//		numSessionCall   int
//		numSessIDCall    int
//		numConnectToCall int
//	}{
//		{
//			args: args{
//				ctx:          context.Background(),
//				primaryIndex: 1,
//				nodes: []nodeArg{
//					{"127.0.0.1", "1"},
//					{"127.0.0.2", "2"},
//					{"127.0.0.3", "3"},
//					{"127.0.0.4", "4"},
//					{"127.0.0.5", "5"},
//					{"127.0.0.6", "6"},
//					{"127.0.0.7", "7"},
//					{"127.0.0.8", "8"},
//					{"127.0.0.9", "9"},
//					{"127.0.0.10", "10"},
//				},
//				primaryPastelID: "2",
//				primarySessID:   "xdcfjc",
//				pastelIDS:       []string{"1", "3", "4", "5", "6", "7", "9"},
//				returnErr:       nil,
//				acceptNodeErr:   nil,
//			},
//			assertion:        assert.NoError,
//			numSessionCall:   10,
//			numSessIDCall:    9,
//			numConnectToCall: 9,
//			want: []string{
//				"1:127.0.0.1",
//				"2:127.0.0.2",
//				"3:127.0.0.3",
//				"4:127.0.0.4",
//				"5:127.0.0.5",
//				"6:127.0.0.6",
//				"7:127.0.0.7",
//				"9:127.0.0.9",
//			},
//		},
//		{
//			args: args{
//				ctx:          context.Background(),
//				primaryIndex: 1,
//				nodes: []nodeArg{
//					{"127.0.0.1", "1"},
//					{"127.0.0.2", "2"},
//					{"127.0.0.3", "3"},
//					{"127.0.0.4", "4"},
//					{"127.0.0.5", "5"},
//					{"127.0.0.6", "6"},
//					{"127.0.0.7", "7"},
//					{"127.0.0.8", "8"},
//					{"127.0.0.9", "9"},
//					{"127.0.0.10", "10"},
//				},
//				primaryPastelID: "2",
//				primarySessID:   "xdxdf",
//				pastelIDS:       []string{"1", "3", "4", "5", "6", "7", "9"},
//				returnErr:       nil,
//				acceptNodeErr:   fmt.Errorf("primary node not accepted"),
//			},
//			assertion:        assert.Error,
//			numSessionCall:   10,
//			numSessIDCall:    9,
//			numConnectToCall: 9,
//			want:             nil,
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("TestTask_meshNodes-%d", i), func(t *testing.T) {
//			nodeClient := nodeTest.NewMockClient(t)
//			nodeClient.
//				ListenOnConnect("", tt.args.returnErr).
//				ListenOnProcessUserdata().
//				ListenOnSessionUserdata(tt.args.returnErr).
//				ListenOnConnectToUserdata(tt.args.returnErr).
//				ListenOnSessIDUserdata(tt.args.primarySessID).
//				ListenOnAcceptedNodesUserdata(tt.args.pastelIDS, tt.args.acceptNodeErr).
//				ListenOnMeshNodes(nil)
//
//			nodes := node.List{}
//			for _, n := range tt.args.nodes {
//				nodes.Add(node.NewNode(nodeClient.Client, n.address, n.pastelID))
//			}
//			service := &UserDataService{
//				config: NewConfig(),
//			}
//			task := &UserDataTask{
//				UserDataService: service,
//			}
//			got, err := task.meshNodes(tt.args.ctx, nodes, tt.args.primaryIndex, &alts.SecInfo{})
//			tt.assertion(t, err)
//			assert.Equal(t, tt.want, pullPastelAddressIDNodes(got))
//
//			nodeClient.AssertAcceptedNodesCallUserdata(1, mock.Anything)
//			nodeClient.AssertSessIDCallUserdata(tt.numSessIDCall)
//			nodeClient.AssertSessionCallUserdata(tt.numSessionCall, mock.Anything, false)
//			nodeClient.AssertConnectToCallUserdata(tt.numConnectToCall, mock.Anything, types.MeshedSuperNode{NodeID: tt.args.primaryPastelID, SessID: tt.args.primarySessID})
//			nodeClient.Client.AssertExpectations(t)
//			nodeClient.Connection.AssertExpectations(t)
//		})
//	}
//}
//
//func TestTask_pastelTopNodes(t *testing.T) {
//	t.Parallel()
//
//	type args struct {
//		ctx       context.Context
//		maxNode   int
//		returnMn  pastel.MasterNodes
//		returnErr error
//	}
//	tests := []struct {
//		args    args
//		want    node.List
//		wantErr bool
//	}{
//		{
//			args: args{
//				ctx:     context.Background(),
//				maxNode: 3,
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4440", ExtKey: "0"},
//					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4441", ExtKey: "1"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4442", ExtKey: "2"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4443", ExtKey: "3"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4444", ExtKey: "4"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "5"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4446", ExtKey: "6"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4447", ExtKey: "7"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4448", ExtKey: "8"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4449", ExtKey: "9"},
//				},
//				returnErr: nil,
//			},
//			want: node.List{
//				newTestNode("127.0.0.1:4440", "0"),
//				newTestNode("127.0.0.1:4441", "1"),
//				newTestNode("127.0.0.1:4442", "2"),
//			},
//			wantErr: false,
//		},
//		{
//			args: args{
//				ctx:     context.Background(),
//				maxNode: 10,
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4440", ExtKey: "0"},
//					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4441", ExtKey: "1"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4442", ExtKey: "2"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4443", ExtKey: "3"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4444", ExtKey: "4"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "5"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4446", ExtKey: "6"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4447", ExtKey: "7"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4448", ExtKey: "8"},
//					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4449", ExtKey: "9"},
//				},
//				returnErr: nil,
//			},
//			want: node.List{
//				newTestNode("127.0.0.1:4440", "0"),
//				newTestNode("127.0.0.1:4441", "1"),
//				newTestNode("127.0.0.1:4442", "2"),
//				newTestNode("127.0.0.1:4443", "3"),
//				newTestNode("127.0.0.1:4444", "4"),
//				newTestNode("127.0.0.1:4445", "5"),
//				newTestNode("127.0.0.1:4446", "6"),
//				newTestNode("127.0.0.1:4447", "7"),
//				newTestNode("127.0.0.1:4448", "8"),
//				newTestNode("127.0.0.1:4449", "9"),
//			},
//			wantErr: false,
//		},
//		{
//			args: args{
//				ctx:       context.Background(),
//				maxNode:   10,
//				returnMn:  nil,
//				returnErr: fmt.Errorf("timeout"),
//			},
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("TestTask_pastelTopNodes-%d", i), func(t *testing.T) {
//			pastelClient := pastelMock.NewMockClient(t)
//			pastelClient.ListenOnMasterNodesTop(tt.args.returnMn, tt.args.returnErr)
//			service := &UserDataService{
//				pastelClient: pastelClient.Client,
//			}
//
//			task := &UserDataTask{
//				UserDataService: service,
//			}
//			got, err := task.pastelTopNodes(tt.args.ctx, tt.args.maxNode)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Task.pastelTopNodes() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Task.pastelTopNodes() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestTask_Error(t *testing.T) {
//	type fields struct {
//		err error
//	}
//	tests := []struct {
//		fields  fields
//		wantErr bool
//	}{
//		{
//			fields: fields{
//				err: fmt.Errorf("error"),
//			},
//			wantErr: true,
//		},
//		{
//			fields: fields{
//				err: nil,
//			},
//			wantErr: false,
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("TestTask_Error-%d", i), func(t *testing.T) {
//			task := &UserDataTask{
//				err: tt.fields.err,
//			}
//			if err := task.Error(); (err != nil) != tt.wantErr {
//				t.Errorf("Task.Error() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestTask_SubscribeProcessResult(t *testing.T) {
//	type fields struct {
//		resultChan chan *userdata.ProcessResult
//	}
//
//	resultChan := make(chan *userdata.ProcessResult)
//
//	tests := []struct {
//		fields fields
//		want   <-chan *userdata.ProcessResult
//	}{
//		{
//			fields: fields{
//				resultChan: resultChan,
//			},
//			want: resultChan,
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("TestTask_SubscribeProcessResult-%d", i), func(t *testing.T) {
//			task := &UserDataTask{
//				resultChan: tt.fields.resultChan,
//			}
//			if got := task.SubscribeProcessResult(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Task.SubscribeProcessResult() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestTask_SubscribeProcessResultGet(t *testing.T) {
//	t.Parallel()
//	type fields struct {
//		resultChanGet chan *userdata.ProcessRequest
//	}
//
//	chanGet := make(chan *userdata.ProcessRequest)
//
//	tests := []struct {
//		fields fields
//		want   <-chan *userdata.ProcessRequest
//	}{
//		{
//			fields: fields{
//				resultChanGet: chanGet,
//			},
//			want: chanGet,
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("TestTask_SubscribeProcessResultGet-%d", i), func(t *testing.T) {
//			task := &UserDataTask{
//				resultChanGet: tt.fields.resultChanGet,
//			}
//			if got := task.SubscribeProcessResultGet(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Task.SubscribeProcessResultGet() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewTask(t *testing.T) {
//	t.Parallel()
//
//	type args struct {
//		service      *UserDataService
//		request      *userdata.ProcessRequest
//		userpastelid string
//	}
//
//	service := &UserDataService{}
//	request := &userdata.ProcessRequest{}
//
//	tests := []struct {
//		args args
//		want *UserDataTask
//	}{
//		{
//			args: args{
//				service:      service,
//				request:      request,
//				userpastelid: "abc",
//			},
//			want: &UserDataTask{
//				UserDataService: service,
//				request:         request,
//				userpastelid:    "abc",
//			},
//		},
//	}
//	for i, tt := range tests {
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			task := NewUserDataTask(tt.args.service, tt.args.request, tt.args.userpastelid)
//			assert.Equal(t, tt.want.UserDataService, task.UserDataService)
//			assert.Equal(t, tt.want.request, task.request)
//			// assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
//		})
//	}
//}
//
//// ReceiveUserdata retrieve the userdata from Metadata layer
//func (nodes *List) ReceiveUserdata(ctx context.Context, pasteluserid string) error {
//	group, _ := errgroup.WithContext(ctx)
//	for _, node := range *nodes {
//		node := node
//		group.Go(func() (err error) {
//			res, err := node.ReceiveUserdata(ctx, pasteluserid)
//			node.ResultGet = res
//			return err
//		})
//	}
//	return group.Wait()
//}
