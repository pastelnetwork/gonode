package healthcheckchallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	healthcheckchallenge "github.com/pastelnetwork/gonode/common/utils/healthcheckchallenge"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix                = "wn-health-check-challenge"
	HealthCheckChallengePort = 9089
)

// GetSummaryStats represents the request for the summary stats method.
type GetSummaryStats struct {
	From       *time.Time
	To         *time.Time
	PastelID   string
	Passphrase string
}

// HCDetailedLogRequest represents the request for the GetDetailedLogs method.
type HCDetailedLogRequest struct {
	ChallengeID string
	PastelID    string
	Passphrase  string
}

// Service represents a service for the SN metrics
type Service struct {
	pastelHandler *mixins.PastelHandler
	client        http.Client
}

// GetSummaryStats returns the summary stats for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetSummaryStats(ctx context.Context, req GetSummaryStats) (healthcheckchallenge.HCSummaryStatsRes, error) {
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return healthcheckchallenge.HCSummaryStatsRes{}, err
	}

	if len(topNodes) < 3 {
		return healthcheckchallenge.HCSummaryStatsRes{}, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return healthcheckchallenge.HCSummaryStatsRes{}, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]healthcheckchallenge.HCSummaryStatsRes)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon healthcheckchallenge.HCSummaryStatsRes
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.fetchSummaryStatsFromNode(ip, req.PastelID, string(signature), req.From, req.To)
			log.WithContext(ctx).Infof("Summary Stats response: %v", data)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node-ip", ip).Error("failed to fetch summary stats from node")

				// increment error counter
				func() {
					mutex.Lock()
					defer mutex.Unlock()

					errorCount++
				}()

				return
			}

			// add results to map and increment success counter
			func() {
				mutex.Lock()
				defer mutex.Unlock()

				results[ip] = data
				successCount++

				hash := data.Hash()
				counts[hash]++
				if counts[hash] > maxCount {
					maxCount = counts[hash]
					mostCommon = data
				}
			}()

		}(strings.Split(node.IPAddress, ":")[0])
	}

	wg.Wait()

	// Proceed only if at least 3 nodes responded successfully
	if successCount < 3 {
		return healthcheckchallenge.HCSummaryStatsRes{}, fmt.Errorf("failed to fetch summary stats from at least 3 nodes, only %d nodes responded successfully", successCount)
	}

	// Log divergent results
	for ip, data := range results {
		if data.Hash() != mostCommon.Hash() {
			log.WithContext(ctx).Infof("Node %s has a divergent result.", ip)
		}
	}

	return mostCommon, nil
}

// fetchSummaryStatsFromNode makes an HTTP GET request to the node's metrics endpoint and returns the metrics.
func (service *Service) fetchSummaryStatsFromNode(addr, pid, passphrase string, from, to *time.Time) (data healthcheckchallenge.HCSummaryStatsRes, err error) {
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/health_check_challenge/summary_stats?pid=%s", addr, HealthCheckChallengePort, pid)
	if from != nil {
		url += "&start=" + from.Format(time.RFC3339)
	}
	if to != nil {
		url += "&end=" + to.Format(time.RFC3339)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", passphrase)

	// Execute the request
	resp, err := service.client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	log.WithContext(context.Background()).Infof("Metrics response: %s", string(body))

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// NewHealthCheckChallengeService returns a new Service instance.
func NewHealthCheckChallengeService(pastelClient pastel.Client) *Service {
	return &Service{
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		client:        http.Client{Timeout: time.Duration(60) * time.Second},
	}
}
