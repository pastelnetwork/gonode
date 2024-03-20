package storagechallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/storagechallenge"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	StorageChallengePort = 9089
)

// GetSummaryStats represents the request for the summary stats method.
type GetSummaryStats struct {
	From       *time.Time
	To         *time.Time
	PastelID   string
	Passphrase string
}

// SCDetailedLogRequest represents the request for the GetDetailedLogs method.
type SCDetailedLogRequest struct {
	ChallengeID string
	Count       int
	PastelID    string
	Passphrase  string
}

// Service represents a service for the SN metrics
type Service struct {
	pastelHandler *mixins.PastelHandler
	client        http.Client
}

// GetSummaryStats returns the summary stats for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetSummaryStats(ctx context.Context, req GetSummaryStats) (storagechallenge.SCSummaryStatsRes, error) {
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return storagechallenge.SCSummaryStatsRes{}, err
	}

	if len(topNodes) < 3 {
		log.WithContext(ctx).Error("failed to get sufficient top nodes")
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("invalid pid/passphrase")
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]storagechallenge.SCSummaryStatsRes)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon storagechallenge.SCSummaryStatsRes
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.fetchSummaryStatsFromNode(ip, req.PastelID, string(signature), req.From, req.To)
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
			log.WithContext(ctx).Infof("Summary Stats response: %v", data)

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
		log.WithContext(ctx).Error("failed to fetch summary stats from at least 3 nodes")
		return storagechallenge.SCSummaryStatsRes{}, fmt.Errorf("failed to fetch summary stats from at least 3 nodes, only %d nodes responded successfully", successCount)
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
func (service *Service) fetchSummaryStatsFromNode(addr, pid, passphrase string, from, to *time.Time) (data storagechallenge.SCSummaryStatsRes, err error) {
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/storage_challenge/summary_stats?pid=%s", addr, StorageChallengePort, pid)
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

// GetDetailedMessageDataList returns the reports for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetDetailedMessageDataList(ctx context.Context, req SCDetailedLogRequest) (messageDataList types.StorageChallengeMessages, err error) {
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return messageDataList, fmt.Errorf("failed to get top nodes: %w", err)
	}

	if len(topNodes) < 3 {
		log.WithContext(ctx).Error("failed to get top nodes")
		return messageDataList, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("invalid pid/passphrase")
		return messageDataList, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]types.StorageChallengeMessages)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon types.StorageChallengeMessages
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.fetchSCDetailedLogs(ip, req.PastelID, string(signature), req.ChallengeID, req.Count)

			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node-ip", ip).Error("failed to fetch detailed-logs from node")

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
		log.WithContext(ctx).Error("failed to fetch detailed-logs from at least 3 nodes")
		return messageDataList, fmt.Errorf("failed to fetch detailed-logs from at least 3 nodes, only %d nodes responded successfully", successCount)
	}

	// Log divergent results
	for ip, data := range results {
		if data.Hash() != mostCommon.Hash() {
			log.WithContext(ctx).Infof("Node %s has a divergent result.", ip)
		}
	}

	return mostCommon, nil
}

// fetchSCMessageLogsData makes an HTTP GET request to the node's detailed_logs endpoint and returns the SCMessageData.
func (service *Service) fetchSCDetailedLogs(addr, pid, passphrase string, challengeID string, count int) (data types.StorageChallengeMessages, err error) {
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/storage_challenge/detailed_logs?pid=%s", addr, StorageChallengePort, pid)
	if challengeID != "" {
		url += "&challenge_id=" + challengeID
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

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return data, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// NewStorageChallengeService returns a new Service instance.
func NewStorageChallengeService(pastelClient pastel.Client) *Service {
	return &Service{
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		client:        http.Client{Timeout: time.Duration(60) * time.Second},
	}
}
