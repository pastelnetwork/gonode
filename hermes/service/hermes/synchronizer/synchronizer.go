package synchronizer

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	//SynchronizationIntervalSec checks after this interval
	SynchronizationIntervalSec = 5
	//SynchronizationTimeoutSec timeout after this
	SynchronizationTimeoutSec = 60
	// MasterNodeSuccessfulStatus represents the MasterNode successful Status
	MasterNodeSuccessfulStatus = "Masternode successfully started"
)

// CheckSynchronized checks the master node status
func (s *Synchronizer) checkSynchronized(ctx context.Context) error {
	st, err := s.PastelClient.MasterNodeStatus(ctx)
	if err != nil {
		return errors.Errorf("getMasterNodeStatus: %w", err)
	}

	if st == nil {
		return errors.New("empty status")
	}

	if st.Status == MasterNodeSuccessfulStatus {
		return nil
	}

	return errors.Errorf("node not synced, status is %s", st.Status)
}

// WaitSynchronization waits for master node to get synced
func (s *Synchronizer) WaitSynchronization(ctx context.Context) error {
	checkTimeout := func(checked chan<- struct{}) {
		time.Sleep(SynchronizationTimeoutSec * time.Second)
		close(checked)
	}

	timeoutCh := make(chan struct{})
	go checkTimeout(timeoutCh)

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(SynchronizationIntervalSec * time.Second):
			err := s.checkSynchronized(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
			} else {
				log.WithContext(ctx).Info("Done for waiting synchronization status")
				return nil
			}
		case <-timeoutCh:
			return errors.New("timeout expired")
		}
	}
}
