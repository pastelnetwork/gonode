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
	"jXYegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc",
}

type Mocker struct {
	pasteldAddrs []string
	ddAddrs      []string
	rqAddrs      []string
	itHelder     *helper.ItHelper
}

func New(pasteldAddrs []string, ddAddrs []string, rqAddrs []string, h *helper.ItHelper) *Mocker {
	return &Mocker{
		pasteldAddrs: pasteldAddrs,
		ddAddrs:      ddAddrs,
		rqAddrs:      rqAddrs,
		itHelder:     h,
	}
}

func (m *Mocker) MockAllRegExpections() error {
	for _, addr := range m.pasteldAddrs {
		if err := m.mockPasteldRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	for _, addr := range m.rqAddrs {
		if err := m.mockRqRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	for _, addr := range m.ddAddrs {
		if err := m.mockDDServerRegExpections(addr); err != nil {
			return fmt.Errorf("server: %s err: %w", addr, err)
		}
	}

	return nil
}

func (m *Mocker) CleanupAll() error {
	addrs := append(m.pasteldAddrs, m.rqAddrs...)
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

func (m *Mocker) mockRqRegExpections(addr string) error {
	if err := m.mockServer([]byte(encodeInfoResp), addr, "encodemetadata",
		[]string{testconst.ArtistPastelID}, 6); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte(encodeResp), addr, "encode",
		[]string{""}, 6); err != nil {
		return fmt.Errorf("failed to mock encode err: %w", err)
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

	if err := m.mockServer([]byte(blockVerboseResponse), addr, "getblock", []string{"160647"}, 10); err != nil {
		return fmt.Errorf("failed to mock masternode top err: %w", err)
	}

	if err := m.mockServer([]byte("160647"), addr, "getblockcount", []string{""}, 10); err != nil {
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
