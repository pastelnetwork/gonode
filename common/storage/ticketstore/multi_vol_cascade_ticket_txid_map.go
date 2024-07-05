package ticketstore

import (
	"time"

	"github.com/pastelnetwork/gonode/common/types"
)

type MultiVolCascadeTicketMapQueries interface {
	InsertMultiVolCascadeTicketTxIDMap(m types.MultiVolCascadeTicketTxIDMap) error
	GetMultiVolCascadeTicketTxIDMap(baseFileID string) (*types.MultiVolCascadeTicketTxIDMap, error)
}

func (s *TicketStore) InsertMultiVolCascadeTicketTxIDMap(m types.MultiVolCascadeTicketTxIDMap) error {
	insertSQL := `
    INSERT INTO multi_vol_cascade_tickets_txid_map (base_file_id, multi_vol_cascade_ticket_txid, created_at, updated_at)
    VALUES (?, ?, ?, ?)`
	stmt, err := s.db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.BaseFileID, m.MultiVolCascadeTicketTxid, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (s *TicketStore) GetMultiVolCascadeTicketTxIDMap(baseFileID string) (*types.MultiVolCascadeTicketTxIDMap, error) {
	selectSQL := `
    SELECT id, base_file_id, multi_vol_cascade_ticket_txid
    FROM multi_vol_cascade_tickets_txid_map
    WHERE base_file_id = ?`
	row := s.db.QueryRow(selectSQL, baseFileID)

	var ticket types.MultiVolCascadeTicketTxIDMap
	err := row.Scan(&ticket.ID, &ticket.BaseFileID, &ticket.MultiVolCascadeTicketTxid)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
