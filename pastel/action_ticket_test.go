package pastel

import (
	"fmt"
	"strings"
	"testing"

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
	inputAPIData := &APISenseTicket{
		DataHash:             []byte{1, 2, 3},
		DDAndFingerprintsIDs: []string{"hello"},
		DDAndFingerprintsIc:  2,
		DDAndFingerprintsMax: 10,
	}

	inputTicket := &ActionTicket{
		Version:    1,
		Caller:     string([]byte{2, 3, 4}),
		BlockNum:   5,
		BlockHash:  string([]byte{6, 7, 8}),
		ActionType: ActionTypeSense,

		APITicketData: inputAPIData,
	}

	_, err := inputTicket.APISenseTicket()
	assert.Nil(t, err)

	_, err = inputTicket.APICascadeTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid action type"))

	inputTicket.ActionType = ActionTypeCascade
	_, err = inputTicket.APICascadeTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid type of api"))
}

func TestAPICascadeTicket(t *testing.T) {
	inputAPIData := &APICascadeTicket{
		DataHash: []byte{1, 2, 3},
	}

	inputTicket := &ActionTicket{
		Version:    1,
		Caller:     string([]byte{2, 3, 4}),
		BlockNum:   5,
		BlockHash:  string([]byte{6, 7, 8}),
		ActionType: ActionTypeCascade,

		APITicketData: inputAPIData,
	}

	_, err := inputTicket.APICascadeTicket()
	assert.Nil(t, err)

	_, err = inputTicket.APISenseTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid action type"))

	inputTicket.ActionType = ActionTypeSense
	_, err = inputTicket.APISenseTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid type of api"))
}
