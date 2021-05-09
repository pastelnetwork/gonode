package fake

import "github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"

func newTicketsListIDMine(tickets models.IDTickets, node *models.MasterNode) models.IDTickets {
	if node == nil {
		return tickets
	}

	return append(tickets, models.IDTicket{
		Height: 19377,
		Ticket: models.IDTicketProp{
			Address:  node.Payee,
			IDType:   "masternode",
			PastelID: node.ExtKey,
			Outpoint: node.Outpoint,
			Type:     "pastelid",
		},
		TXID: "53e6f0a01d9a7a309552dfe4a80003f50db585b3f200723825f5587a4a06c02e",
	})
}
