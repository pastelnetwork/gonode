package collection

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/pastelnetwork/gonode/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCollectionTicketService(t *testing.T) {
	s, tmpFileName := prepareCollectionService(t)
	defer os.Remove(tmpFileName)

	actTicket := &pastel.ActTicket{}
	faker.FakeData(actTicket)
	actTicket.TXID = "test-act-tx-id"
	actTicket.ActTicketData.RegTXID = "test-reg-tx-id"

	collectionRegTicket := &pastel.CollectionRegTicket{}
	faker.FakeData(collectionRegTicket)
	collectionRegTicket.TXID = "test-reg-tx-id"
	collectionRegTicket.CollectionRegTicketData.CollectionTicketData.CollectionItemCopyCount = 10
	collectionRegTicket.CollectionRegTicketData.CollectionTicketData.MaxCollectionEntries = 10

	appTicketData := pastel.AppTicket{
		MaxPermittedOpenNSFWScore:                      0.5,
		MinimumSimilarityScoreToFirstEntryInCollection: 0.3,
	}

	appTicketBytes, err := json.Marshal(appTicketData)
	require.NoError(t, err)

	appTicket := base64.RawStdEncoding.EncodeToString(appTicketBytes)
	collectionRegTicket.CollectionRegTicketData.CollectionTicketData.AppTicket = appTicket

	tests := []struct {
		testcase string
		expected bool
		setup    func(t *testing.T)
		expect   func(*testing.T, error)
	}{
		{
			testcase: "when no collection act tickets returned by pastel, the service should end gracefully",
			expected: false,
			setup: func(t *testing.T) {
				pMock := pastelMock.NewMockClient(t)

				pMock.ListenOnCollectionActivationTicketsFromBlockHeight(pastel.ActTickets{}, nil)

				s.pastelClient = pMock
			},
			expect: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			testcase: "when act tickets returned by pastel, with collection max and current no of collection entries are equal," +
				" it should be stored in the DB, with collection state as finalized",
			expected: false,
			setup: func(t *testing.T) {
				pMock := pastelMock.NewMockClient(t)

				pMock.ListenOnCollectionActivationTicketsFromBlockHeight(pastel.ActTickets{*actTicket}, nil)
				pMock.ListenOnCollectionRegTicket(*collectionRegTicket, nil)

				s.pastelClient = pMock
			},
			expect: func(t *testing.T, err error) {
				require.NoError(t, err)
				collectionTicket, err := s.store.GetCollection(context.Background(), "test-act-tx-id")
				require.NoError(t, err)
				require.Equal(t, domain.InProcessCollectionState, collectionTicket.CollectionState)
				require.Equal(t, actTicket.Height, collectionTicket.CollectionTicketActivationBlockHeight)
				require.Equal(t, collectionTicket.MinimumSimilarityScoreToFirstEntryInCollection, 0.3)
				require.Equal(t, collectionTicket.MaxPermittedOpenNSFWScore, 0.5)
			},
		},
	}

	for _, test := range tests {
		test := test // add this if there's subtest (t.Run)
		t.Run(test.testcase, func(t *testing.T) {
			test.setup(t)

			err := s.parseCollectionTickets(context.Background())

			test.expect(t, err)
		})

	}
}

func prepareCollectionService(t *testing.T) (*collectionService, string) {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	store, err := store.NewSQLiteStore(tmpfile.Name())
	assert.Nil(t, err)

	tmpfile.Close()

	s := &collectionService{
		store: store,
	}

	return s, tmpfile.Name()
}
