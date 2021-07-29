package dupedetection

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
)

type client struct {
	config *Config
}

// DDServiceResult is JSON output structure of dupe detection service
type DDServiceResult struct {
	DDSystemVersion             string              `json:"dupe_detection_system_version"`
	ImageHash                   string              `json:"hash_of_candidate_image_file"`
	IsLikelyDupe                float64             `json:"is_likely_dupe"`
	OverallAverageRarenessScore float64             `json:"overall_average_rareness_score"`
	IsRareOnInternet            float64             `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage     int                 `json:"matches_found_on_first_page"`
	NumberOfResultPages         int                 `json:"number_of_pages_of_results"`
	FirstMatchURL               string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore               float64             `json:"open_nsfw_score"`
	AlternateNSFWScores         AlternateNSFWScores `json:"alternative_nsfw_scores"`
	FingerPrints                string              `json:"image_fingerprint_of_candidate_image_file"`
}

// Generate implements dupedection.Client.Generate
func (client *client) Generate(ctx context.Context, path string) (*DupeDetectionResult, error) {

	// Copy iamge file to dupe detection service input directory
	outputBase, outputPath, err := client.copyImageToInputDir(path)
	if err != nil {
		return nil, errors.Errorf("failed to copy image to dupe detection input directory: %w", err)
	}

	// Waiting for JSON output from dupe detection service
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil, errors.New("Timeout")
		case <-ticker.C:
			files, err := ioutil.ReadDir(client.config.OutputDir)
			if err != nil {
				return nil, errors.Errorf("failed to scan dupe detection output directory: %v", err)
			}
			for _, file := range files {
				if file.Name() == outputBase {
					jsonFile, err := os.Open(outputPath)
					if err != nil {
						return nil, errors.Errorf("failed to open JSON output: %v", err)
					}
					defer jsonFile.Close()

					byteValue, err := ioutil.ReadAll(jsonFile)
					if err != nil {
						return nil, errors.Errorf("failed to read JSON output: %v", err)
					}

					ddresult := DDServiceResult{}
					err = json.Unmarshal(byteValue, &ddresult)
					if err != nil {
						return nil, errors.Errorf("failed to parse JSON output: %v", err)
					}

					result := &DupeDetectionResult{
						PastelRarenessScore:   ddresult.OverallAverageRarenessScore,
						InternetRarenessScore: ddresult.IsRareOnInternet,
						OpenNSFWScore:         ddresult.OpenNSFWScore,
						AlternateNSFWScores:   ddresult.AlternateNSFWScores,
						FingerPrints:          ddresult.FingerPrints,
					}
					return result, nil
				}
			}
		}
	}
}

func (client *client) copyImageToInputDir(path string) (string, string, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	ext := filepath.Ext(path)

	// Generate random name
	inputName, err := random.String(10, random.Base62Chars)
	if err != nil {
		return "", "", err
	}

	inputPath := filepath.Join(client.config.InputDir, inputName+ext)

	err = ioutil.WriteFile(inputPath, input, os.ModePerm)
	if err != nil {
		return "", "", err
	}

	outputBase := inputName + ".json"
	outputPath := client.config.OutputDir + outputBase
	return outputBase, outputPath, nil
}

// NewClient return a new Client instance
func NewClient(config *Config) Client {
	return &client{
		config: config,
	}
}
