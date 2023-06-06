package mixins

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// GetPastelIDfromMNConfig gets the pastelID from the MN config API
func GetPastelIDfromMNConfig(ctx context.Context, pastelClient pastel.Client, confKey string) (string, error) {
	log.WithContext(ctx).Info("Reading PastelID from MN Config")
	// Define the backoff policy
	backoffPolicy := backoff.NewConstantBackOff(5 * time.Second)

	// Set max retry attempts
	maxRetryBackoff := backoff.WithMaxRetries(backoffPolicy, 5)
	var extKey string

	// Operation to retry
	operation := func() error {
		mnStatus, err := pastelClient.MasterNodeStatus(ctx)
		if err != nil {
			log.Warnf("Error getting master-node status: %v", err) // Updated log message
			return err
		}

		if mnStatus.ExtKey == "" {
			return errors.New("extKey is empty")
		}

		if mnStatus.ExtKey != confKey {
			log.Warnf("Warning! pastel IDs do not match - ID in config: %s - ID from cnode API( this will be used ): %s\n", confKey, extKey)
		}
		extKey = mnStatus.ExtKey // Save the extKey

		return nil // We found a matching extKey, hence no retry
	}

	// Perform the operation with retry
	err := backoff.Retry(operation, maxRetryBackoff)
	if err != nil {
		log.Warnf("Unable to get pastelID from mn config API - error after retries: %v", err) // Updated log message
		return "", err
	}

	return extKey, nil
}
