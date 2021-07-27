package pastel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeArtTicket(t *testing.T) {
	inputAppData := AppTicket{
		BlockNum:              10,
		PreviewHash:           []byte{1},
		Thumbnail1Hash:        []byte{2},
		Thumbnail2Hash:        []byte{3},
		DataHash:              []byte{4},
		Fingerprints:          []byte{5},
		FingerprintsHash:      []byte{6},
		FingerprintsSignature: []byte{7},

		RQIDs: []string{"9", "10"},
	}

	inputTicket := ArtTicket{
		Version:       1,
		Author:        []byte{2, 3, 4},
		BlockNum:      5,
		BlockHash:     []byte{6, 7, 8},
		Copies:        9,
		Royalty:       10,
		Green:         "11",
		AppTicketData: inputAppData,
	}

	encoded, err := EncodeArtTicket(&inputTicket)
	assert.Nil(t, err)
	outputTicket, err := DecodeArtTicket(encoded)
	outputAppData := outputTicket.AppTicketData
	assert.Nil(t, err)
	fmt.Println(string(encoded))
	assert.Equal(t, inputTicket.Version, outputTicket.Version)
	assert.Equal(t, inputTicket.Author, outputTicket.Author)
	assert.Equal(t, inputTicket.BlockNum, outputTicket.BlockNum)
	assert.Equal(t, inputTicket.BlockHash, outputTicket.BlockHash)
	assert.Equal(t, inputTicket.Copies, outputTicket.Copies)
	assert.Equal(t, inputTicket.Royalty, outputTicket.Royalty)
	assert.Equal(t, inputTicket.Green, outputTicket.Green)

	assert.Equal(t, inputAppData.BlockNum, outputAppData.BlockNum)
	assert.Equal(t, inputAppData.PreviewHash, outputAppData.PreviewHash)
	assert.Equal(t, inputAppData.Thumbnail1Hash, outputAppData.Thumbnail1Hash)
	assert.Equal(t, inputAppData.DataHash, outputAppData.DataHash)
	assert.Equal(t, inputAppData.Fingerprints, outputAppData.Fingerprints)
}
