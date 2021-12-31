package pastel

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeActionTicket(t *testing.T) {
	inputApiData := &ApiSenseTicket{
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

		ApiTicketData: inputApiData,
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
	outputApiData, ok := outputTicket.ApiTicketData.(*ApiSenseTicket)
	assert.True(t, ok)

	assert.Equal(t, inputApiData.DataHash, outputApiData.DataHash)
	assert.Equal(t, inputApiData.DDAndFingerprintsIc, outputApiData.DDAndFingerprintsIc)
	assert.Equal(t, inputApiData.DDAndFingerprintsMax, outputApiData.DDAndFingerprintsMax)
}

func TestApiSenseTicket(t *testing.T) {
	inputApiData := &ApiSenseTicket{
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

		ApiTicketData: inputApiData,
	}

	_, err := inputTicket.ApiSenseTicket()
	assert.Nil(t, err)

	_, err = inputTicket.ApiCascadeTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid action type"))

	inputTicket.ActionType = ActionTypeCascade
	_, err = inputTicket.ApiCascadeTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid type of api"))
}

func TestApiCascadeTicket(t *testing.T) {
	inputApiData := &ApiCascadeTicket{
		DataHash: []byte{1, 2, 3},
	}

	inputTicket := &ActionTicket{
		Version:    1,
		Caller:     string([]byte{2, 3, 4}),
		BlockNum:   5,
		BlockHash:  string([]byte{6, 7, 8}),
		ActionType: ActionTypeCascade,

		ApiTicketData: inputApiData,
	}

	_, err := inputTicket.ApiCascadeTicket()
	assert.Nil(t, err)

	_, err = inputTicket.ApiSenseTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid action type"))

	inputTicket.ActionType = ActionTypeSense
	_, err = inputTicket.ApiSenseTicket()
	assert.True(t, strings.Contains(err.Error(), "invalid type of api"))
}
