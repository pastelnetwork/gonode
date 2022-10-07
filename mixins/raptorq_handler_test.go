package mixins

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/files"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestGenRQIdentifiersFiles(t *testing.T) {
	type args struct {
		connectErr error
	}

	testCases := map[string]struct {
		args              args
		wantErr           error
		encodeInfoReturns *rqnode.EncodeInfo
	}{
		"success": {
			args:    args{},
			wantErr: nil,
			encodeInfoReturns: &rqnode.EncodeInfo{
				SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
					"test-file": {
						ID:                uuid.New().String(),
						SymbolIdentifiers: []string{"test-s1, test-s2"},
						BlockHash:         "test-block-hash",
						PastelID:          "test-pastel-id",
					},
				},
				EncoderParam: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
			},
		},
		"error": {
			args: args{
				connectErr: errors.New("test"),
			},
			encodeInfoReturns: &rqnode.EncodeInfo{
				SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
					"test-file": {
						ID:                uuid.New().String(),
						SymbolIdentifiers: []string{"test-s1, test-s2"},
						BlockHash:         "test-block-hash",
						PastelID:          "test-pastel-id",
					},
				},
				EncoderParam: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil)
			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(tc.encodeInfoReturns, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr)

			fsMock := storageMock.NewMockFileStorage()

			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)
			storage := files.NewStorage(fsMock)
			fsMock.ListenOnOpen(fileMock, nil)
			file := files.NewFile(storage, "test")

			ph := NewPastelHandler(pastelClientMock)

			h := NewRQHandler(rqClientMock, ph, "", "", 1, 1)
			err := h.GenRQIdentifiersFiles(context.Background(), file, "", "", "")
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
