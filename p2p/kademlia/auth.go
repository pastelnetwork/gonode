package kademlia

import (
	"bytes"
	"context"
	"encoding/binary"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	timestampMarginDuration = 5 * time.Second
)

type AuthHelper struct {
	pastelClient pastel.Client
	secInfo      *alts.SecInfo

	expiredDuration time.Duration
	authMtx         sync.Mutex
	timestamp       time.Time
	signature       []byte

	masterNodes                []pastel.MasterNode
	masterNodesMtx             sync.Mutex
	lastMasterNodesRefresh     time.Time
	masterNodesRefreshDuration time.Duration
}

// NewAuthHelper returns a peer authHelper
func NewAuthHelper(pastelClient pastel.Client, secInfo *alts.SecInfo) *AuthHelper {
	return &AuthHelper{
		pastelClient:               pastelClient,
		secInfo:                    secInfo,
		expiredDuration:            5 * time.Minute,
		masterNodes:                []pastel.MasterNode{},
		masterNodesRefreshDuration: 5 * time.Minute,
	}
}

// GenAuthInfo generatue auth info of peer to show his finger print
func (ath *AuthHelper) GenAuthInfo(ctx context.Context) (*auth.PeerAuthInfo, error) {
	ath.authMtx.Lock()
	defer ath.authMtx.Unlock()
	// if need to refresh current timestamp
	if time.Now().After(ath.timestamp.Add(ath.expiredDuration + timestampMarginDuration)) {
		if err := ath.unsafeRefreshAuthInfo(ctx); err != nil {
			return nil, errors.Errorf("refresh auth: %w", err)
		}
	}

	return &auth.PeerAuthInfo{
		PastelID:  ath.secInfo.PastelID,
		Timestamp: ath.timestamp,
		Signature: ath.signature,
	}, nil
}

// VerifyPeer verify auth info of other side
func (ath *AuthHelper) VerifyPeer(ctx context.Context, authInfo *auth.PeerAuthInfo) error {
	// verify timestamp - if too old
	if time.Now().Add(-ath.expiredDuration).After(authInfo.Timestamp) {
		return errors.New("invalid timestamp - too old")
	}

	// verify timestamp - if too high
	if authInfo.Timestamp.After(time.Now().Add(timestampMarginDuration)) {
		return errors.New("invalid timestamp - far from now")
	}

	// check if peer from list
	if !ath.isPeerInMasterNodes(ctx, authInfo) {
		return errors.New("peer not in master nodes")
	}

	// verify signature
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, authInfo.Timestamp.Unix())

	if ok, err := ath.pastelClient.Verify(ctx, buf.Bytes(), string(authInfo.Signature), ath.secInfo.PastelID, ath.secInfo.Algorithm); err != nil || !ok {
		if err == nil {
			err = errors.New("signature not match")
		}
		return errors.Errorf("signature verify: %w", err)
	}

	return nil
}

func (ath *AuthHelper) isPeerInMasterNodes(ctx context.Context, authInfo *auth.PeerAuthInfo) bool {
	// try to refresh master nodes if need
	ath.tryRefreshMasternodeList(ctx)

	ath.masterNodesMtx.Lock()
	defer ath.masterNodesMtx.Unlock()

	if len(ath.masterNodes) == 0 {
		return false
	}

	// check if in master nodes
	for _, node := range ath.masterNodes {
		// if PastelID match, verify IP
		if node.ExtKey == authInfo.PastelID {

			// extract IP address
			addr := strings.Split(node.ExtAddress, ":")
			if len(addr) != 2 {
				return false
			}
			ip := addr[0]

			// empty ip
			if ip == "" {
				return false
			}

			// Check if peer address has same IP
			if strings.HasSuffix(authInfo.Address, ip) {
				return true
			}

			return false
		}
	}

	return false
}

func (ath *AuthHelper) unsafeRefreshAuthInfo(ctx context.Context) error {
	newTimestamp := time.Now()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, newTimestamp.Unix())
	newSignature, err := ath.pastelClient.Sign(ctx, buf.Bytes(), ath.secInfo.PastelID, ath.secInfo.PassPhrase, ath.secInfo.Algorithm)
	if err != nil {
		return err
	}

	ath.timestamp = newTimestamp
	ath.signature = []byte(newSignature)

	return nil
}

func (ath *AuthHelper) tryRefreshMasternodeList(ctx context.Context) {
	ath.masterNodesMtx.Lock()
	defer ath.masterNodesMtx.Unlock()
	if len(ath.masterNodes) == 0 || time.Now().After(ath.lastMasterNodesRefresh.Add(ath.masterNodesRefreshDuration)) {
		if err := ath.unsafeRefreshMasternodeList(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Error("Update master node list failed")
			return
		}
		ath.lastMasterNodesRefresh = time.Now()
	}
}
func (ath *AuthHelper) unsafeRefreshMasternodeList(ctx context.Context) error {
	masterNodes, err := ath.pastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	ath.masterNodes = masterNodes
	return nil
}
