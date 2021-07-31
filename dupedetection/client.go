package dupedetection

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
)

type client struct {
	config *Config
}

type dupeDetectionResult struct {
	DupeDetectionSystemVer  string              `json:"dupe_detection_system_version"`
	ImageHash               string              `json:"hash_of_candidate_image_file,omitempty"`
	IsLikelyDupe            float64             `json:"is_likely_dupe,omitempty"`
	PastelRarenessScore     float64             `json:"overall_average_rareness_score"`
	InternetRarenessScore   float64             `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage int                 `json:"matches_found_on_first_page"`
	NumberOfResultPages     int                 `json:"number_of_pages_of_results"`
	FirstMatchURL           string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore           float64             `json:"open_nsfw_score"`
	AlternateNSFWScores     AlternateNSFWScores `json:"alternative_nsfw_scores"`
	ImageHashes             ImageHashes         `json:"image_hashes"`
	FingerPrints            string              `json:"image_fingerprint_of_candidate_image_file"`
}

// Generate implements dupedection.Client.Generate
func (client *client) Generate(ctx context.Context, path string) (*DupeDetection, error) {

	// Copy image file to dupe detection service input directory
	outputBase, outputPath, err := client.copyImageToInputDir(path)
	if err != nil {
		return nil, errors.Errorf("failed to copy image to dupe detection input directory: %w", err)
	}

	// Collect result from dupe detection output directory
	result, err := client.collectOutput(ctx, outputBase, outputPath)
	if err != nil {
		err = errors.Errorf("failed to collect output from dupe detection service: %w", err)
	}

	return result, err
}

func (client *client) copyImageToInputDir(inputPath string) (string, string, error) {
	input, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return "", "", err
	}

	ext := filepath.Ext(inputPath)

	// Generate random name
	inputName, err := random.String(10, random.Base62Chars)
	if err != nil {
		return "", "", err
	}

	inputPath = filepath.Join(client.config.InputDir, inputName+ext)

	err = ioutil.WriteFile(inputPath, input, os.ModePerm)
	if err != nil {
		return "", "", err
	}

	outputBase := inputName + ".json"
	outputPath := filepath.Join(client.config.OutputDir, outputBase)
	return outputBase, outputPath, nil
}

func (client *client) collectOutput(ctx context.Context, baseName, path string) (*DupeDetection, error) {
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
				if file.Name() == baseName {
					result, err := parseOutput(path)

					if err != nil {
						log.WithContext(ctx).WithError(err).Debug("Failed to collect output from dupe detection service")
						continue
					}

					return result, nil
				}
			}
		}
	}
}

func parseOutput(path string) (*DupeDetection, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Errorf("failed to open dupe detection JSON output: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Errorf("failed to read dupe detection JSON output: %v", err)
	}

	ddResult := dupeDetectionResult{}
	err = json.Unmarshal(byteValue, &ddResult)
	if err != nil {
		return nil, errors.Errorf("failed to parse dupe detection service JSON output: %v", err)
	}

	result := DupeDetection{
		DupeDetectionSystemVer:  ddResult.DupeDetectionSystemVer,
		PastelRarenessScore:     ddResult.PastelRarenessScore,
		InternetRarenessScore:   ddResult.InternetRarenessScore,
		MatchesFoundOnFirstPage: ddResult.MatchesFoundOnFirstPage,
		NumberOfResultPages:     ddResult.NumberOfResultPages,
		FirstMatchURL:           ddResult.FirstMatchURL,
		OpenNSFWScore:           ddResult.OpenNSFWScore,
		AlternateNSFWScores:     ddResult.AlternateNSFWScores,
		ImageHashes:             ddResult.ImageHashes,
	}
	err = json.Unmarshal([]byte(ddResult.FingerPrints), &result.FingerPrints)
	if err != nil {
		return nil, errors.Errorf("failed to parse fingerprints from dupe detection service JSON output: %v", err)
	}

	return &result, nil
}

// NewClient return a new Client instance
func NewClient(config *Config) Client {
	return &client{
		config: config,
	}
}
