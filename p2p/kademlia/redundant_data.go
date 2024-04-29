package kademlia

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

func (s *DHT) startDisabledKeysCleanupWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("disabled keys cleanup worker started")

	for {
		select {
		case <-time.After(defaultCleanupInterval):
			s.cleanupDisabledKeys(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing disabled keys cleanup worker")
			return nil
		}
	}
}

func (s *DHT) cleanupDisabledKeys(ctx context.Context) error {
	if s.metaStore == nil {
		return nil
	}

	from := time.Now().UTC().Add(-1 * defaultDisabledKeyExpirationInterval)
	disabledKeys, err := s.metaStore.GetDisabledKeys(from)
	if err != nil {
		return errors.Errorf("get disabled keys: %w", err)
	}

	for i := 0; i < len(disabledKeys); i++ {
		dec, err := hex.DecodeString(disabledKeys[i].Key)
		if err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("decode disabled key failed")
			continue
		}
		s.metaStore.Delete(ctx, dec)
	}

	return nil
}

func (s *DHT) startCleanupRedundantDataWorker(ctx context.Context) {
	log.P2P().WithContext(ctx).Info("redundant data cleanup worker started")

	for {
		select {
		case <-time.After(defaultRedundantDataCleanupInterval):
			s.cleanupRedundantDataWorker(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing disabled keys cleanup worker")
			return
		}
	}
}

func (s *DHT) cleanupRedundantDataWorker(ctx context.Context) {
	from := time.Now().AddDate(-5, 0, 0) // 5 years ago

	log.P2P().WithContext(ctx).WithField("from", from).Info("getting all possible replication keys past five years")
	to := time.Now().UTC()
	replicationKeys := s.store.GetKeysForReplication(ctx, from, to)

	ignores := s.ignorelist.ToNodeList()
	self := &Node{ID: s.ht.self.ID, IP: s.externalIP, Port: s.ht.self.Port}
	self.SetHashedID()

	closestContactsMap := make(map[string][][]byte)

	for i := 0; i < len(replicationKeys); i++ {
		decKey, _ := hex.DecodeString(replicationKeys[i].Key)
		nodes := s.ht.closestContactsWithInlcudingNode(Alpha, decKey, ignores, self)
		closestContactsMap[replicationKeys[i].Key] = nodes.NodeIDs()
	}

	insertKeys := make([]domain.DelKey, 0)
	removeKeys := make([]domain.DelKey, 0)
	for key, closestContacts := range closestContactsMap {
		if len(closestContacts) < Alpha {
			log.P2P().WithContext(ctx).WithField("key", key).WithField("closest contacts", closestContacts).Info("not enough contacts to replicate")
			continue
		}

		found := false
		for _, contact := range closestContacts {
			if bytes.Equal(contact, self.ID) {
				found = true
			}
		}

		delKey := domain.DelKey{
			Key:       key,
			CreatedAt: time.Now(),
			Count:     1,
		}

		if !found {
			insertKeys = append(insertKeys, delKey)
		} else {
			removeKeys = append(removeKeys, delKey)
		}
	}

	if len(insertKeys) > 0 {
		if err := s.metaStore.BatchInsertDelKeys(ctx, insertKeys); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("insert keys failed")
			return
		}

		log.P2P().WithContext(ctx).WithField("count-del-keys", len(insertKeys)).Info("insert del keys success")
	} else {
		log.P2P().WithContext(ctx).Info("No redundant key found to be stored in the storage")
	}

	if len(removeKeys) > 0 {
		if err := s.metaStore.BatchDeleteDelKeys(ctx, removeKeys); err != nil {
			log.WithContext(ctx).WithError(err).Error("batch delete del-keys failed")
			return
		}
	}

}

func (s *DHT) startDeleteDataWorker(ctx context.Context) {
	log.P2P().WithContext(ctx).Info("start delete data worker")

	for {
		select {
		case <-time.After(defaultDeleteDataInterval):
			s.deleteRedundantData(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing delete data worker")
			return
		}
	}
}

func (s *DHT) deleteRedundantData(ctx context.Context) {
	const batchSize = 100

	delKeys, err := s.metaStore.GetAllToDelKeys(delKeysCountThreshold)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("get all to delete keys failed")
		return
	}

	for len(delKeys) > 0 {
		// Check the available disk space
		isLow, err := CheckDiskSpace()
		if err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("check disk space failed")
			break
		}

		if !isLow {
			// Disk space is sufficient, stop deletion
			break
		}

		// Determine the size of the next batch
		batchEnd := batchSize
		if len(delKeys) < batchSize {
			batchEnd = len(delKeys)
		}

		// Prepare the batch for deletion
		keysToDelete := make([]string, 0, batchEnd)
		for _, delKey := range delKeys[:batchEnd] {
			keysToDelete = append(keysToDelete, delKey.Key)
		}

		// Perform the deletion
		if err := s.store.BatchDeleteRecords(keysToDelete); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("batch delete records failed")
			break
		}

		// Update the remaining keys to be deleted
		delKeys = delKeys[batchEnd:]
	}
}

// CheckDiskSpace checks if the available space on the disk is less than 50 GB
func CheckDiskSpace() (bool, error) {
	// Define the command and its arguments
	cmd := exec.Command("bash", "-c", "df -BG --output=source,fstype,avail | egrep 'ext4|xfs' | sort -rn -k3 | awk 'NR==1 {print $3}'")

	// Execute the command
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to execute command: %w", err)
	}

	// Process the output
	output := strings.TrimSpace(out.String())
	if len(output) < 2 {
		return false, errors.New("invalid output from disk space check")
	}

	// Convert the available space to an integer
	availSpace, err := strconv.Atoi(output[:len(output)-1])
	if err != nil {
		return false, fmt.Errorf("failed to parse available space: %w", err)
	}

	// Check if the available space is less than 50 GB
	return availSpace < lowSpaceThreshold, nil
}
