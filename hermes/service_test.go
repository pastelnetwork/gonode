package hermes

import (
	"context"
	"errors"
	"fmt"
	"testing"

	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestCleanupInactiveTickets(t *testing.T) {
	type args struct {
		regTicketsErr        error
		actionTicketsErr     error
		regTicketsReturns    pastel.RegTickets
		actionTicketsReturns pastel.ActionTicketDatas
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				regTicketsErr:        nil,
				actionTicketsErr:     nil,
				regTicketsReturns:    getRegTickets(),
				actionTicketsReturns: getActionTickets(),
			},
			wantErr: nil,
		},

		"reg-tickets-err": {
			args: args{
				regTicketsErr:        errors.New("test"),
				actionTicketsErr:     nil,
				regTicketsReturns:    getRegTickets(),
				actionTicketsReturns: getActionTickets(),
			},
			wantErr: errors.New("test"),
		},

		"action-tickets-err": {
			args: args{
				regTicketsErr:        nil,
				actionTicketsErr:     errors.New("test"),
				regTicketsReturns:    getRegTickets(),
				actionTicketsReturns: getActionTickets(),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnRegTicketsFromBlockHeight(tc.args.regTicketsReturns, 1, tc.args.regTicketsErr).
				ListenOnActionTicketsFromBlockHeight(tc.args.actionTicketsReturns, 1, tc.args.actionTicketsErr)

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnDelete(nil)

			service := NewService(NewConfig(), pastelClientMock, p2pClient)

			err := service.cleanupInactiveTickets(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func getActionTickets() (actionTickets pastel.ActionTicketDatas) {
	cascadeTxid := "b2fe52d207beaf107ad248612e2ad3580562aa2e0e1a7a567b89997c2fb382cd"
	ct := pastel.ActionRegTicket{
		Height: 125,
		TXID:   cascadeTxid,
		ActionTicketData: pastel.ActionTicketData{
			CalledAt:   125,
			Type:       "action-reg",
			ActionType: "cascade",
			Label:      "label",
			StorageFee: 5,
			ActionTicketData: pastel.ActionTicket{
				Caller:     "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
				ActionType: "cascade",
				BlockNum:   1,
				BlockHash:  "",
				APITicketData: pastel.APICascadeTicket{
					RQIDs: []string{
						"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z",
					},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220,
						47, 235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
				},
			},
		},
	}
	actionByte, err := pastel.EncodeActionTicket(&ct.ActionTicketData.ActionTicketData)
	if err != nil {
		panic(err)
	}
	ct.ActionTicketData.ActionTicket = actionByte

	actionTickets = append(actionTickets, ct)

	senseTxid := "c2fe52d207beaf107ad248612e2ad3580562aa2e0e1a7a567b89997c2fb382po"
	st := pastel.ActionRegTicket{
		Height: 120,
		TXID:   senseTxid,
		ActionTicketData: pastel.ActionTicketData{
			CalledAt:   120,
			Type:       "action-reg",
			ActionType: "sense",
			Label:      "label",
			StorageFee: 5,
			ActionTicketData: pastel.ActionTicket{
				Caller:     "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
				ActionType: "sense",
				BlockNum:   1,
				BlockHash:  "",
				APITicketData: pastel.APISenseTicket{
					DDAndFingerprintsIDs: []string{"ddfpfilehash"},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220,
						47, 235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
				},
			},
		},
	}
	actionByte, err = pastel.EncodeActionTicket(&st.ActionTicketData.ActionTicketData)
	if err != nil {
		panic(err)
	}
	st.ActionTicketData.ActionTicket = actionByte

	return append(actionTickets, st)
}

func getRegTickets() (regTickets pastel.RegTickets) {
	txid := "b4b1fc370983c7409ec58fcd079136f04efe1e1c363f4cd8f4aff8986a91ef09"
	ticket := pastel.RegTicket{
		TXID:   txid,
		Height: 100,
		RegTicketData: pastel.RegTicketData{
			TotalCopies:   5,
			CreatorHeight: 1,
			StorageFee:    5,
			NFTTicketData: pastel.NFTTicket{
				BlockNum: 1,
				Copies:   2,
				AppTicketData: pastel.AppTicket{
					CreatorName:   "Alan Majchrowicz",
					NFTTitle:      "nft-1",
					NFTSeriesName: "nft-series",
					RQIDs: []string{
						"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z",
					},
					Thumbnail1Hash: []byte{
						100, 96, 227, 164, 165, 224, 135, 208, 183, 238, 17, 48, 154, 166, 40,
						20, 0, 215, 158, 151, 198, 123, 115, 153, 162, 81, 118, 200, 22, 153, 174, 130,
					},
					Thumbnail2Hash: []byte{249, 23, 138, 33, 126, 119, 199, 251, 112, 210, 244, 236, 144, 50,
						149, 172, 64, 148, 230, 63, 44, 221, 144, 63, 173, 131, 97, 25, 227, 105, 181, 127},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220,
						47, 235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
					DDAndFingerprintsIDs: []string{"ddfpfilehash"},
				},
			},
		},
	}

	artTicketBytes, _ := pastel.EncodeNFTTicket(&ticket.RegTicketData.NFTTicketData)

	ticket.RegTicketData.NFTTicket = artTicketBytes
	regTickets = append(regTickets, ticket)

	a := pastel.RegTicket{
		TXID:   "reg_a",
		Height: 80,
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					Thumbnail1Hash: []byte{
						100, 96, 227, 164, 165, 224, 135, 208, 183, 238, 17, 48, 154, 166, 40,
						20, 0, 215, 158, 151, 198, 123, 115, 153, 162, 81, 118, 200, 22, 153, 174, 130,
					},
					Thumbnail2Hash: []byte{249, 23, 138, 33, 126, 119, 199, 251, 112, 210, 244, 236, 144, 50,
						149, 172, 64, 148, 230, 63, 44, 221, 144, 63, 173, 131, 97, 25, 227, 105, 181, 127},
					CreatorName: "Alan Majchrowicz",
					NFTTitle:    "Lake Superior Sky III",
					TotalCopies: 10,
					CreatorWrittenStatement: `Lakes Why settle for blank walls, when you can
					 transform them into stunning vista points. Explore Lake Superiorfrom imaginative
					  scenic abstracts to sublime beach landscapes captured on camera.`,
					NFTKeywordSet:        "Michigan,Midwest,Peninsula,Great Lakes,Lakeview",
					NFTSeriesName:        "Science Art Lake",
					DDAndFingerprintsIDs: []string{"ddfpfilehash"},
				},
			},
			TotalCopies: 10,
		},
	}

	artTicketBytes, _ = pastel.EncodeNFTTicket(&a.RegTicketData.NFTTicketData)
	a.RegTicketData.NFTTicket = artTicketBytes

	regTickets = append(regTickets, a)

	b := pastel.RegTicket{
		Height: 90,
		TXID:   "reg_b",
		RegTicketData: pastel.RegTicketData{
			TotalCopies:   5,
			CreatorHeight: 1,
			StorageFee:    5,
			NFTTicketData: pastel.NFTTicket{

				BlockNum: 1,
				Copies:   2,
				AppTicketData: pastel.AppTicket{
					Thumbnail1Hash: []byte{
						100, 96, 227, 164, 165, 224, 135, 208, 183, 238, 17, 48, 154, 166, 40,
						20, 0, 215, 158, 151, 198, 123, 115, 153, 162, 81, 118, 200, 22, 153, 174, 130,
					},
					Thumbnail2Hash: []byte{249, 23, 138, 33, 126, 119, 199, 251, 112, 210, 244, 236, 144, 50,
						149, 172, 64, 148, 230, 63, 44, 221, 144, 63, 173, 131, 97, 25, 227, 105, 181, 127},
					CreatorName:   "Mona Lisa",
					NFTTitle:      "Art-b",
					NFTSeriesName: "nft-series",
					RQIDs: []string{
						"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z",
					},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220,
						47, 235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
					DDAndFingerprintsIDs: []string{"ddfpfilehash"},
				},
			},
		},
	}

	artTicketBytes, _ = pastel.EncodeNFTTicket(&b.RegTicketData.NFTTicketData)
	b.RegTicketData.NFTTicket = artTicketBytes
	regTickets = append(regTickets, b)

	c := pastel.RegTicket{
		TXID:   "reg_c",
		Height: 75,
		RegTicketData: pastel.RegTicketData{
			TotalCopies:   5,
			CreatorHeight: 1,
			StorageFee:    5,
			NFTTicketData: pastel.NFTTicket{
				BlockNum: 1,
				Copies:   2,
				AppTicketData: pastel.AppTicket{
					Thumbnail1Hash: []byte{
						100, 96, 227, 164, 165, 224, 135, 208, 183, 238, 17, 48, 154, 166, 40,
						20, 0, 215, 158, 151, 198, 123, 115, 153, 162, 81, 118, 200, 22, 153, 174, 130,
					},
					Thumbnail2Hash: []byte{249, 23, 138, 33, 126, 119, 199, 251, 112, 210, 244, 236, 144, 50,
						149, 172, 64, 148, 230, 63, 44, 221, 144, 63, 173, 131, 97, 25, 227, 105, 181, 127},
					CreatorName:   "Alan Border",
					NFTTitle:      "nft-c",
					NFTSeriesName: "nft-series",
					RQIDs: []string{
						"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z",
					},
					DataHash: []byte{237, 203, 202, 207, 43, 36, 172, 173, 251, 169, 72, 216, 216, 220,
						47, 235, 33, 171, 187, 188, 199, 189, 250, 43, 39, 154, 14, 144, 51, 135, 12, 68},
					DDAndFingerprintsIDs: []string{"ddfpfilehash"},
				},
			},
		},
	}

	artTicketBytes, _ = pastel.EncodeNFTTicket(&c.RegTicketData.NFTTicketData)
	c.RegTicketData.NFTTicket = artTicketBytes
	regTickets = append(regTickets, c)

	return regTickets
}
