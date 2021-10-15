package externaldupedetection

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"

	"github.com/pastelnetwork/gonode/dupedetection/ddclient"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/ddclient/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/tj/assert"
)

func TestTaskSignAndSendEDDTicket(t *testing.T) {
	type args struct {
		task        *Task
		signErr     error
		sendArtErr  error
		signReturns []byte
		primary     bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				primary: true,
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				signErr: errors.New("test"),
				primary: true,
			},
			wantErr: errors.New("test"),
		},
		"primary-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				sendArtErr: errors.New("test"),
				signErr:    nil,
				primary:    false,
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.args.signReturns, tc.args.signErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			clientMock := test.NewMockClient(t)
			clientMock.ListenOnExternalDupeDetection_SendEDDTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnExternalDupeDetection()

			tc.args.task.nodeClient = clientMock
			tc.args.task.connectedTo = &Node{client: clientMock}
			err := tc.args.task.connectedTo.connect(context.Background())
			assert.NoError(t, err)

			err = tc.args.task.signAndSendDDTicket(context.Background(), tc.args.primary)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskRegisterArt(t *testing.T) {
	type args struct {
		task     *Task
		regErr   error
		regRetID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
				regErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnRegisterArtTicket(tc.args.regRetID, tc.args.regErr).
				ListenOnRegisterExDDTicket(tc.args.regRetID, tc.args.regErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			id, err := tc.args.task.registerExternalDupeDetection(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.args.regRetID, id)
			}
		})
	}
}

func TestTaskGenFingerprintsData(t *testing.T) {
	type args struct {
		task    *Task
		fileErr error
		genErr  error
		genResp *ddclient.DupeDetection
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				genErr:  nil,
				fileErr: errors.New("test"),
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			file := artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			fingerprints := []float32{12.3, 34.4}
			//fgBytes, err := json.Marshal(fingerprints)
			//assert.Nil(t, err)

			tc.args.genResp.Fingerprints = fingerprints
			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock
			_, _, err := tc.args.task.genFingerprintsData(context.Background(), file)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskPastelNodesByExtKey(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				masterNodesErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"node-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				nodeID:         "B",
				masterNodesErr: nil,
			},
			wantErr: errors.New("not found"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			_, err := tc.args.task.pastelNodeByExtKey(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskCompareRQSymbolID(t *testing.T) {
	type args struct {
		task        *Task
		connectErr  error
		fileErr     error
		assignRQIDS bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"decompress-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				connectErr:  nil,
				fileErr:     nil,
				assignRQIDS: true,
			},
			wantErr: errors.New("failed to decompress"),
		},
		"conn-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				connectErr:  errors.New("test"),
				fileErr:     nil,
				assignRQIDS: true,
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				fileErr:     errors.New("test"),
				connectErr:  nil,
				assignRQIDS: true,
			},
			wantErr: errors.New("read image"),
		},
		"rqids-len-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				fileErr:     nil,
				connectErr:  nil,
				assignRQIDS: false,
			},
			wantErr: errors.New("no symbols identifiers file"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr)

			tc.args.task.Service.rqClient = rqClientMock
			tc.args.task.rqClient = rqClientMock

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			tc.args.task.Artwork = artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			if tc.args.assignRQIDS {
				rqFile := rq.SymbolIDFile{ID: "A"}
				bytes, err := json.Marshal(rqFile)
				assert.NoError(t, err)

				tc.args.task.RQIDS["A"] = bytes
			}

			err := tc.args.task.compareRQSymbolID(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskStoreRaptorQSymbols(t *testing.T) {
	type args struct {
		task       *Task
		encodeErr  error
		connectErr error
		fileErr    error
		storeErr   error
		encodeResp *rqnode.Encode
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    errors.New("test"),
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"conn-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				encodeErr:  nil,
				connectErr: errors.New("test"),
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"encode-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				encodeErr:  errors.New("test"),
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"store-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
					RQIDS:  make(map[string][]byte),
				},
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   errors.New("test"),
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqFile := rq.SymbolIDFile{ID: "A"}
			bytes, err := json.Marshal(rqFile)
			assert.Nil(t, err)
			tc.args.task.RQIDS["A"] = bytes
			tc.args.encodeResp.Symbols = map[string][]byte{"A": bytes}

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr).
				ListenOnEncode(tc.args.encodeResp, tc.args.encodeErr)

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnStore("", tc.args.storeErr)
			tc.args.task.Service.p2pClient = p2pClient
			tc.args.task.p2pClient = p2pClient

			tc.args.task.Service.rqClient = rqClientMock
			tc.args.task.rqClient = rqClientMock

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			tc.args.task.Artwork = artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err = tc.args.task.storeRaptorQSymbols(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskStoreFingerprints(t *testing.T) {
	type args struct {
		task     *Task
		storeErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				storeErr: nil,
			},
			wantErr: nil,
		},

		"store-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				storeErr: errors.New("test"),
			},
			wantErr: errors.New("failed to store fingerprints"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnStore("", tc.args.storeErr)
			tc.args.task.Service.p2pClient = p2pClient
			tc.args.task.p2pClient = p2pClient

			err := tc.args.task.storeFingerprints(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskVerifyPeersSignature(t *testing.T) {
	type args struct {
		task      *Task
		verifyErr error
		verifyRet bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					peersEDDTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: true,
				verifyErr: nil,
			},
			wantErr: nil,
		},
		"verify-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					peersEDDTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: true,
				verifyErr: errors.New("test"),
			},
			wantErr: errors.New("verify signature"),
		},
		"verify-failure": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					peersEDDTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: false,
				verifyErr: nil,
			},
			wantErr: errors.New("mistmatch"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(tc.args.verifyRet, tc.args.verifyErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			err := tc.args.task.verifyPeersSingature(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskWaitConfirmation(t *testing.T) {
	type args struct {
		task             *Task
		txid             string
		timeout          time.Duration
		interval         time.Duration
		minConfirmations int64
		ctxTimeout       time.Duration
	}

	testCases := map[string]struct {
		args    args
		wantErr error
		retRes  *pastel.GetRawTransactionVerbose1Result
		retErr  error
	}{
		"timeout": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				interval:   200 * time.Millisecond,
				timeout:    100 * time.Millisecond,
				ctxTimeout: 500 * time.Millisecond,
			},
			wantErr: errors.New("timeout"),
		},
		"min-confirmations-timeout": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				minConfirmations: 2,
				timeout:          200 * time.Millisecond,
				interval:         100 * time.Millisecond,
				ctxTimeout:       500 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: errors.New("timeout"),
		},
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				minConfirmations: 1,
				timeout:          100 * time.Millisecond,
				interval:         50 * time.Millisecond,
				ctxTimeout:       500 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: nil,
		},
		"ctx-done-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ExDDTicket{},
				},
				minConfirmations: 1,
				timeout:          100 * time.Millisecond,
				interval:         50 * time.Millisecond,
				ctxTimeout:       10 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: errors.New("context"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), tc.args.ctxTimeout)
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetRawTransactionVerbose1(tc.retRes, tc.retErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			err := <-tc.args.task.waitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.timeout, tc.args.interval)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			cancel()
		})
	}

}

func TestTaskProbeImage(t *testing.T) {
	type args struct {
		task    *Task
		fileErr error
		genErr  error
		genResp *ddclient.DupeDetection
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:                    task.New(StatusConnected),
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:                    task.New(StatusImageProbed),
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("required status"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: &ddclient.DupeDetection{},
				task: &Task{
					Task: task.New(StatusConnected),
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.ExDDTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersEDDTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			file := artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			fingerprints := []float32{12.3, 34.4}
			tc.args.genResp.Fingerprints = fingerprints
			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			_, err := tc.args.task.ProbeImage(context.Background(), file)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskGetRegistrationFee(t *testing.T) {
	type args struct {
		task   *Task
		retFee int64
		retErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusImageUploaded),
					Ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: nil,
		},

		"status-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusConnected),
					Ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: errors.New("require status"),
		},
		"fee-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusImageUploaded),
					Ticket: &pastel.ExDDTicket{},
				},
				retErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetRegisterExDDFee(tc.args.retFee, tc.args.retErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)
			eddTicketBytes, err := pastel.EncodeExDDTicket(tc.args.task.Ticket)
			assert.NoError(t, err)

			_, err = tc.args.task.GetRegistrationFee(context.Background(), eddTicketBytes, tc.args.task.creatorPastelID,
				[]byte{}, "", "", map[string][]byte{}, []byte{})
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskSessionNode(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusPrimaryMode),
					Ticket: &pastel.ExDDTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusConnected),
					Ticket: &pastel.ExDDTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: errors.New("status"),
		},
		"pastel-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusPrimaryMode),
					Ticket: &pastel.ExDDTicket{},
				},
				masterNodesErr: errors.New("test"),
				nodeID:         "A",
			},
			wantErr: errors.New("failed to get node"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			err := tc.args.task.SessionNode(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskAddPeerArtTicketSignature(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
		acceptedNodeID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					peersEDDTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                  task.New(StatusRegistrationFeeCalculated),
					Ticket:                &pastel.ExDDTicket{},
					allSignaturesReceived: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				task: &Task{
					peersEDDTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                  task.New(StatusConnected),
					Ticket:                &pastel.ExDDTicket{},
					allSignaturesReceived: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: errors.New("status"),
		},
		"no-node-err": {
			args: args{
				task: &Task{
					peersEDDTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                  task.New(StatusRegistrationFeeCalculated),
					Ticket:                &pastel.ExDDTicket{},
					allSignaturesReceived: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("accepted"),
		},
		"success-close-sign-chn": {
			args: args{
				task: &Task{
					peersEDDTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					allSignaturesReceived: make(chan struct{}),
					Task:                  task.New(StatusRegistrationFeeCalculated),
					Ticket:                &pastel.ExDDTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: tc.args.nodeID})
			tc.args.task.accepted = Nodes{&Node{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID}}

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.args.task.peersEDDTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := tc.args.task.AddPeerArtTicketSignature(tc.args.nodeID, []byte{})
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskUploadImage(t *testing.T) {
	type args struct {
		task *Task
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusImageProbed),
					Ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: nil,
		},
		"failure-status": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusConnected),
					Ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: errors.New("status"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			stg := artwork.NewStorage(fs.NewFileStorage(os.TempDir()))

			tc.args.task.Storage = stg
			file, err := newTestImageFile(stg)
			assert.NoError(t, err)
			tc.args.task.Artwork = file

			err = tc.args.task.UploadImage(ctx, file)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newTestImageFile(stg *artwork.Storage) (*artwork.File, error) {
	imgFile := stg.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create storage file: %w", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	if err := png.Encode(f, img); err != nil {
		return nil, err
	}

	return imgFile, nil
}
