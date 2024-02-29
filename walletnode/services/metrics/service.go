package metrics

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
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/metrics"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix   = "wn-metrics"
	metricsPort = 9089
)

// GetSummaryStats represents the request for the summary stats method.
type GetSummaryStats struct {
	From       *time.Time
	To         *time.Time
	PastelID   string
	Passphrase string
}

// SHReportRequest represents the request for the GetDetailedLogs method.
type SHReportRequest struct {
	Count      int
	EventID    string
	PastelID   string
	Passphrase string
}

// Service represents a service for the SN metrics
type Service struct {
	pastelHandler *mixins.PastelHandler
	client        http.Client
}

// GetSummaryStats returns the summary stats for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetSummaryStats(ctx context.Context, req GetSummaryStats) (metrics.Metrics, error) {
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return metrics.Metrics{}, err
	}

	if len(topNodes) < 3 {
		return metrics.Metrics{}, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return metrics.Metrics{}, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]metrics.Metrics)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon metrics.Metrics
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.fetchMetricsFromNode(ip, req.PastelID, string(signature), req.From, req.To)
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
		return metrics.Metrics{}, fmt.Errorf("failed to fetch summary stats from at least 3 nodes, only %d nodes responded successfully", successCount)
	}

	// Log divergent results
	for ip, data := range results {
		if data.Hash() != mostCommon.Hash() {
			log.WithContext(ctx).Infof("Node %s has a divergent result.", ip)
		}
	}

	return mostCommon, nil
}

// fetchMetricsFromNode makes an HTTP GET request to the node's metrics endpoint and returns the metrics.
func (service *Service) fetchMetricsFromNode(addr, pid, passphrase string, from, to *time.Time) (data metrics.Metrics, err error) {
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/metrics?pid=%s", addr, metricsPort, pid)
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

	log.WithContext(context.Background()).Infof("Summary-Stats response: %s", string(body))

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

// GetDetailedLogs returns the reports for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetDetailedLogs(ctx context.Context, req SHReportRequest) (report types.SelfHealingReports, err error) {
	report = types.SelfHealingReports{}
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return report, fmt.Errorf("failed to get top nodes: %w", err)
	}

	if len(topNodes) < 3 {
		return report, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return report, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]types.SelfHealingReports)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon types.SelfHealingReports
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.fetchSHChallengesFromNode(ip, req.PastelID, string(signature), req.Count, req.EventID)

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
		return report, fmt.Errorf("failed to fetch detailed-logs from at least 3 nodes, only %d nodes responded successfully", successCount)
	}

	// Log divergent results
	for ip, data := range results {
		if data.Hash() != mostCommon.Hash() {
			log.WithContext(ctx).Infof("Node %s has a divergent result.", ip)
		}
	}

	return mostCommon, nil
}

// fetchMetricsFromNode makes an HTTP GET request to the node's metrics endpoint and returns the metrics.
func (service *Service) fetchSHChallengesFromNode(addr, pid, passphrase string, count int, challengeID string) (data types.SelfHealingReports, err error) {
	data = types.SelfHealingReports{}
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/sh_challenge?pid=%s", addr, metricsPort, pid)
	if count != 0 {
		url += "&count=" + fmt.Sprintf("%d", count)
	} else if challengeID != "" {
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
		return data, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// NewMetricsService returns a new Service instance.
func NewMetricsService(pastelClient pastel.Client) *Service {
	return &Service{
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		client:        http.Client{Timeout: time.Duration(60) * time.Second},
	}
}
