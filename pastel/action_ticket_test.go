package pastel

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncodeActionTicket(t *testing.T) {
	inputAPIData := &APISenseTicket{
		DataHash:             []byte{1, 2, 3},
		DDAndFingerprintsIDs: []string{"hello"},
		DDAndFingerprintsIc:  2,
		DDAndFingerprintsMax: 10,
	}

	inputTicket := ActionTicket{
		Version:    1,
		Caller:     string([]byte{2, 3, 4}),
		BlockNum:   5,
		BlockHash:  string([]byte{6, 7, 8}),
		ActionType: ActionTypeSense,

		APITicketData: inputAPIData,
	}

	encoded, err := EncodeActionTicket(&inputTicket)
	assert.Nil(t, err)
	outputTicket, err := DecodeActionTicket(encoded)

	assert.Nil(t, err)
	fmt.Println(string(encoded))
	assert.Equal(t, inputTicket.Version, outputTicket.Version)
	assert.Equal(t, inputTicket.Caller, outputTicket.Caller)
	assert.Equal(t, inputTicket.BlockNum, outputTicket.BlockNum)
	assert.Equal(t, inputTicket.BlockHash, outputTicket.BlockHash)
	assert.Equal(t, inputTicket.ActionType, outputTicket.ActionType)
	outputAPIData, ok := outputTicket.APITicketData.(*APISenseTicket)
	assert.True(t, ok)

	assert.Equal(t, inputAPIData.DataHash, outputAPIData.DataHash)
	assert.Equal(t, inputAPIData.DDAndFingerprintsIc, outputAPIData.DDAndFingerprintsIc)
	assert.Equal(t, inputAPIData.DDAndFingerprintsMax, outputAPIData.DDAndFingerprintsMax)
}

func TestAPISenseTicket(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in      *ActionTicket
		wantErr error
	}{
		"Success": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeSense,

				APITicketData: &APISenseTicket{
					DataHash:             []byte{1, 2, 3},
					DDAndFingerprintsIDs: []string{"hello"},
					DDAndFingerprintsIc:  2,
					DDAndFingerprintsMax: 10,
				},
			},
			wantErr: nil,
		},
		"invalid action type": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeCascade,

				APITicketData: &APISenseTicket{
					DataHash:             []byte{1, 2, 3},
					DDAndFingerprintsIDs: []string{"hello"},
					DDAndFingerprintsIc:  2,
					DDAndFingerprintsMax: 10,
				},
			},
			wantErr: errors.New("invalid action type"),
		},
		"invalid type of api": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeSense,

				APITicketData: &APICascadeTicket{},
			},
			wantErr: errors.New("invalid type of api"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.in.APISenseTicket()
			if tc.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			}
		})
	}
}

func TestAPICascadeTicket(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in      *ActionTicket
		wantErr error
	}{
		"Success": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeCascade,

				APITicketData: &APICascadeTicket{
					DataHash: []byte{1, 2, 3},
				},
			},
			wantErr: nil,
		},
		"invalid action type": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeSense,

				APITicketData: &APICascadeTicket{
					DataHash: []byte{1, 2, 3},
				},
			},
			wantErr: errors.New("invalid action type"),
		},
		"invalid type of api": {
			in: &ActionTicket{
				Version:    1,
				Caller:     string([]byte{2, 3, 4}),
				BlockNum:   5,
				BlockHash:  string([]byte{6, 7, 8}),
				ActionType: ActionTypeCascade,

				APITicketData: &APISenseTicket{},
			},
			wantErr: errors.New("invalid type of api"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.in.APICascadeTicket()
			if tc.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			}
		})
	}
}
