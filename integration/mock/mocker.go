package mock

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"
	"github.com/pastelnetwork/gonode/integration/helper"
)

var validPastelIDs = []string{
	"jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN",
	"jXZFvrCSNQGfrRA8RaXpzm2TEV38hQiyesNZg7joEjBNjXGieM8ZTyhKYxkBrqdtTpYcACRKQjMLQrmKNpfrL8",
	"jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU",
	"jABegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
	"jXa7Ypp5fueoTboLYxAbH2ebaGLiQBRDF2u2PFAMYYGRjSMsEzn5sgSgyAy5xYpEfcMUqaUY6xZ8Qcxnn1pWAf",
	"jXXQD8CjECcKknAUf5NHcCua8aqHhbq4vP5EyE6kCDvD6EWpnbjaUe8VH9ZVpvj8n85hkMtDZyPqrwJXPwYUBi",
	"jXYqgELnZQUSC2D5ZKPzVxaJwDLKC2NAzTzvEkR5EXU3od7xDJH9JgsMWvU6FoYQaBVP489A1y7qw756m8Er1S",
}

type Mocker struct {
	pasteldAddrs []string
	ddAddrs      []string
	snAddrs      []string
	itHelder     *helper.ItHelper
}

func New(pasteldAddrs []string, ddAddrs []string, snAddrs []string, h *helper.ItHelper) *Mocker {
	return &Mocker{
		pasteldAddrs: pasteldAddrs,
		ddAddrs:      ddAddrs,
		snAddrs:      snAddrs,
		itHelder:     h,
	}
}

func (m *Mocker) MockAllRegExpections() error {
	for _, addr := range m.pasteldAddrs {
		if err := m.mockPasteldRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	/*for _, addr := range m.rqAddrs {
		if err := m.mockRqRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}
	*/
	for _, addr := range m.ddAddrs {
		if err := m.mockDDServerRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	for _, addr := range m.snAddrs {
		if err := m.mockSNFilesExpectations(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	return nil
}

func (m *Mocker) CleanupAll() error {
	addrs := m.pasteldAddrs
	addrs = append(addrs, m.ddAddrs...)
	for _, addr := range addrs {
		if err := m.cleanup(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	return nil
}

func (m *Mocker) mockDDServerRegExpections(addr string) error {
	resp, err := getDDServerResponse()
	if err != nil {
		return err
	}

	if err := m.mockServer(resp, addr, "rareness",
		[]string{testconst.ArtistPastelID}, 6); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	return nil
}

func (m *Mocker) mockPasteldRegExpections(addr string) error {

	if err := m.mockServer([]byte(masterNodesTopResp), addr, "masternode", []string{"top"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(masterNodeExtraResp), addr, "masternode", []string{"list", "extra"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(networkStorageResp), addr, "storagefee", []string{"getnetworkfee"}, 6); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	for _, id := range validPastelIDs {
		if err := m.mockServer([]byte(ticketsFindID), addr, "tickets", []string{"find", "id", id}, 6); err != nil {
			return fmt.Errorf("failed to mock masternode top err: %w", err)
		}
	}

	if err := m.mockServer([]byte(ticketsFindIDArtistResp), addr, "tickets",
		[]string{"find", "id", testconst.ArtistPastelID}, 3); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(`[]`), addr, "tickets",
		[]string{"findbylabel", "nft"}, 3); err != nil {
		return fmt.Errorf("failed to mock findbylabel err: %w", err)
	}

	if err := m.mockServer([]byte(`[]`), addr, "tickets",
		[]string{"findbylabel", "action"}, 3); err != nil {
		return fmt.Errorf("failed to mock findbylabel err: %w", err)
	}

	if err := m.mockServer([]byte(blockVerboseResponse), addr, "getblock", []string{"160647"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte("160647"), addr, "getblockcount", []string{""}, 100); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte("18792"), addr, "z_getbalance", []string{testconst.RegSpendableAddress}, 3); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(signResp), addr, "pastelid", []string{"sign"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(verifyResp), addr, "pastelid", []string{"verify"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte("opid-0ded54d5-f87a-479d-bcc7-e77866e51409"), addr, "z_sendmanywithchangetosender",
		[]string{testconst.RegSpendableAddress}, 5); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(opStatusResp), addr, "z_getoperationstatus",
		[]string{`[opid-0ded54d5-f87a-479d-bcc7-e77866e51409]`}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(storageFeeResp), addr, "tickets", []string{"tools", "gettotalstoragefee"}, 3); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(getRawTxResp), addr, "getrawtransaction", []string{}, 5); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(getNFTRegisterResp), addr, "tickets", []string{"register", "nft"}, 1); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(getActRegisterResp), addr, "tickets", []string{"register", "act"}, 1); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	// sense
	if err := m.mockServer([]byte(actionFeeResp), addr, "storagefee", []string{"getactionfees"}, 6); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(getActionRegisterResp), addr, "tickets", []string{"register", "action"}, 1); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(getActRegisterResp), addr, "tickets", []string{"activate", "action"}, 1); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(ticketsFindIDArtistResp), addr, "tickets",
		[]string{"find", "action-act"}, 3); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	// download
	if err := m.mockServer([]byte(ticketOwnershipResp), addr, "tickets",
		[]string{"tools", "validateownership", "b4b1fc370983c7409ec58fcd079136f04efe1e1c363f4cd8f4aff8986a91ef09",
			testconst.ArtistPastelID, "passphrase"}, 3); err != nil {

		return fmt.Errorf("failed to mock ticket ownership err: %w", err)
	}

	if err := m.mockServer([]byte(cascadeticketOwnershipResp), addr, "tickets",
		[]string{"tools", "validateownership", testconst.TestCascadeRegTXID,
			testconst.ArtistPastelID, "passphrase"}, 3); err != nil {

		return fmt.Errorf("failed to mock ticket ownership err: %w", err)
	}

	for key, val := range getRegTickets() {
		if err := m.mockServer(val, addr, "tickets", []string{"get",
			key}, 20); err != nil {

			return fmt.Errorf("failed to mock tickets get ownership err: %w", err)
		}
	}

	if err := m.mockServer(getTradeTickets(), addr, "tickets", []string{"list",
		"trade", "available"}, 3); err != nil {

		return fmt.Errorf("failed to mock tickets get ownership err: %w", err)
	}

	//search
	if err := m.mockServer(getSearchActTicketsResponse(), addr, "tickets", []string{"list",
		"act"}, 5); err != nil {

		return fmt.Errorf("failed to mock tickets act err: %w", err)
	}

	return nil
}

func (m *Mocker) mockSNFilesExpectations(addr string) error {
	//send each ticket's AppTicketData.DDAndFingerprintsIDs to storage

	storeDDAndFpFile(ddAndFpFile, m.itHelder, addr)

	return nil
}

func (m *Mocker) mockServer(payload []byte, addr string, command string, params []string, times int) error {
	query := make(map[string]string)
	query["method"] = command
	query["params"] = strings.Join(params, ",")
	query["count"] = fmt.Sprint(times)

	rsp, status, err := m.itHelder.RequestRaw(helper.HttpPost, payload, fmt.Sprintf("%s/%s", addr, "register"), query)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New(string(rsp))
	}

	return nil
}

func (m *Mocker) cleanup(addr string) error {
	rsp, status, err := m.itHelder.RequestRaw(helper.HttpPost, []byte{}, fmt.Sprintf("%s/%s", addr, "cleanup"), nil)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New(string(rsp))
	}

	return nil
}

func (m *Mocker) mockSCPasteldRegExpections(addr string) error {

	if err := m.mockServer([]byte(masterNodesTopRespSC), addr, "masternode", []string{"top"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(masterNodeExtraRespSC), addr, "masternode", []string{"list", "extra"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	for _, id := range validPastelIDs {
		if err := m.mockServer([]byte(ticketsFindID), addr, "tickets", []string{"find", "id", id}, 6); err != nil {
			return fmt.Errorf("failed to mock masternode top err: %w", err)
		}
	}

	if err := m.mockServer([]byte(ticketsFindIDArtistResp), addr, "tickets",
		[]string{"find", "id", testconst.ArtistPastelID}, 3); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(`[]`), addr, "tickets",
		[]string{"findbylabel", "nft"}, 3); err != nil {
		return fmt.Errorf("failed to mock findbylabel err: %w", err)
	}

	if err := m.mockServer([]byte(blockVerboseResponseSC), addr, "getblock", []string{"160647"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte("160647"), addr, "getblockcount", []string{""}, 100); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(signResp), addr, "pastelid", []string{"sign"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(verifyResp), addr, "pastelid", []string{"verify"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	tix, tlist := getSCRegTickets()

	for key, val := range tix {
		if err := m.mockServer(val, addr, "tickets", []string{"get",
			key}, 20); err != nil {

			return fmt.Errorf("failed to mock tickets get ownership err: %w", err)
		}
	}

	if err := m.mockServer(tlist, addr, "tickets", []string{"list",
		"nft"}, 10); err != nil {

		return fmt.Errorf("failed to mock tickets act err: %w", err)
	}

	return nil
}

func (m *Mocker) MockSCRegExpections() error {
	for _, addr := range m.pasteldAddrs {
		if err := m.mockSCPasteldRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	/*for _, addr := range m.snAddrs {
		storeRQIDFile([]byte("rqidfile"), m.itHelder, addr)
	}
	*/
	return nil
}
