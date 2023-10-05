package mixins

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// GetPastelIDfromMNConfig gets the pastelID from the MN config API
func GetPastelIDfromMNConfig(ctx context.Context, pastelClient pastel.Client, confKey string) (string, error) {
	log.WithContext(ctx).Info("Reading PastelID from MN Config")

	const maxRetries = 5
	const delayBetweenRetries = 5 * time.Second

	var extKey string

	for i := 0; i < maxRetries; i++ {
		mnStatus, err := pastelClient.MasterNodeStatus(ctx)
		if err != nil {
			log.Warnf("Error getting master-node status: %v", err) // Updated log message
		} else {
			if mnStatus.ExtKey == "" {
				log.Warn("extKey is empty")
			} else {
				if mnStatus.ExtKey != confKey {
					log.Warnf("Warning! pastel IDs do not match - ID in config: %s - ID from cnode API( this will be used ): %s\n", confKey, extKey)
				}
				extKey = mnStatus.ExtKey // Save the extKey
				return extKey, nil       // Successfully found a matching extKey, return it
			}
		}

		// If we haven't returned by this point, sleep before retrying
		if i < maxRetries-1 { // Don't sleep after the last attempt
			time.Sleep(delayBetweenRetries)
		}
	}

	log.Warnf("Unable to get pastelID from mn config API - error after retries.")
	return "", fmt.Errorf("failed to get pastelID after %d retries", maxRetries)
}
