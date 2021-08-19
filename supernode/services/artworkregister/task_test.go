package artworkregister

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/dupedetection"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/tj/assert"
)

func TestTaskSignAndSendArtTicket(t *testing.T) {
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
			clientMock.ListenOnSendArtTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterArtwork()

			tc.args.task.nodeClient = clientMock
			tc.args.task.connectedTo = &Node{client: clientMock}
			err := tc.args.task.connectedTo.connect(context.Background())
			assert.Nil(t, err)

			err = tc.args.task.signAndSendArtTicket(context.Background(), tc.args.primary)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
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
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
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
				ListenOnRegisterNFTTicket(tc.args.regRetID, tc.args.regErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			id, err := tc.args.task.registerArt(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
		genResp *dupedetection.DupeDetection
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: &dupedetection.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				genErr:  nil,
				fileErr: errors.New("test"),
				genResp: &dupedetection.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: &dupedetection.DupeDetection{},
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
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

			fingerprints := []float64{12.3, 34.4}
			fgBytes, err := json.Marshal(fingerprints)
			assert.Nil(t, err)

			tc.args.genResp.Fingerprints = string(fgBytes)
			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnGenerate(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock
			_, _, err = tc.args.task.genFingerprintsData(context.Background(), file)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
				assert.Nil(t, err)

				tc.args.task.RQIDS["A"] = bytes
			}

			err := tc.args.task.compareRQSymbolID(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskStoreThumbnails(t *testing.T) {
	type args struct {
		task     *Task
		storeErr error
		fileErr  error
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
					Ticket: &pastel.NFTTicket{},
				},
				fileErr:  nil,
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
					Ticket: &pastel.NFTTicket{},
				},
				fileErr:  nil,
				storeErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				fileErr:  errors.New("test"),
				storeErr: nil,
			},
			wantErr: errors.New("test"),
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

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)

			tc.args.task.SmallThumbnail = artwork.NewFile(storage, "test-small")
			tc.args.task.MediumThumbnail = artwork.NewFile(storage, "test-medium")
			tc.args.task.PreviewThumbnail = artwork.NewFile(storage, "test-preview")

			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err := tc.args.task.storeThumbnails(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket: &pastel.NFTTicket{},
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
					Ticket: &pastel.NFTTicket{},
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
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
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
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
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
