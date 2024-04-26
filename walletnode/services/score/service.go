package score

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	ScoreChallengesPort = 9089
)

// Service represents a service for the SN metrics
type Service struct {
	pastelHandler *mixins.PastelHandler
	client        http.Client
}

type ChallengesScoreRequest struct {
	PastelID   string
	Passphrase string
}

// GetAggregatedChallengeScore returns the reports for the given PastelID, fetching them concurrently from all nodes.
func (service *Service) GetAggregatedChallengeScore(ctx context.Context, req ChallengesScoreRequest) (scoreList types.AggregatedScoreList, err error) {
	topNodes, err := service.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get top nodes")
		return scoreList, fmt.Errorf("failed to get top nodes: %w", err)
	}

	if len(topNodes) < 3 {
		log.WithContext(ctx).Error("failed to get top nodes")
		return scoreList, fmt.Errorf("no top nodes found")
	}

	signature, err := service.pastelHandler.PastelClient.Sign(ctx, []byte(req.PastelID), req.PastelID, req.Passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("invalid pid/passphrase")
		return scoreList, fmt.Errorf("invalid pid/passphrase: %w", err)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex // Mutex to protect access to the results map
	results := make(map[string]types.AggregatedScoreList)
	errorCount := 0
	successCount := 0

	counts := make(map[string]int) // To count occurrences of each unique result.
	var mostCommon types.AggregatedScoreList
	maxCount := 0

	for _, node := range topNodes {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			data, err := service.getAggregatedChallengesScores(ip, req.PastelID, string(signature))

			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node-ip", ip).Error("failed to fetch challenges scores from node")

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
		log.WithContext(ctx).Error("failed to fetch scores from at least 3 nodes")
		return scoreList, fmt.Errorf("failed to fetch scores from at least 3 nodes, only %d nodes responded successfully", successCount)
	}

	// Log divergent results
	for ip, data := range results {
		if data.Hash() != mostCommon.Hash() {
			log.WithContext(ctx).Infof("Node %s has a divergent result.", ip)
		}
	}

	return mostCommon, nil
}

// GetAggregatedChallengesScores makes an HTTP GET request to the challenges score endpoint and returns the aggregated SC and HC challenge score.
func (service *Service) getAggregatedChallengesScores(addr, pid, passphrase string) (data types.AggregatedScoreList, err error) {
	// Construct the URL with query parameters
	url := fmt.Sprintf("http://%s:%d/nodes/challenges_score?pid=%s", addr, ScoreChallengesPort, pid)

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

// NewScoreService returns a new Service instance.
func NewScoreService(pastelClient pastel.Client) *Service {
	return &Service{
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		client:        http.Client{Timeout: time.Duration(60) * time.Second},
	}
}
