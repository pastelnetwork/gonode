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

// Generate implements dupedection.Client.Generate
func (client *client) Generate(ctx context.Context, path string) (*DupeDetection, error) {

	// Copy image file to dupe detection service input directory
	outputPath, err := client.copyImageToInputDir(path)
	if err != nil {
		return nil, errors.Errorf("failed to copy image to dupe detection input directory: %w", err)
	}

	// Collect result from dupe detection output directory
	result, err := client.collectOutput(ctx, outputPath)
	if err != nil {
		err = errors.Errorf("failed to collect output from dupe detection service: %w", err)
	}

	return result, err
}

func (client *client) copyImageToInputDir(inputPath string) (string, error) {
	input, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(inputPath)

	// Generate random name
	inputName, err := random.String(10, random.Base62Chars)
	if err != nil {
		return "", err
	}

	inputPath = filepath.Join(client.config.InputDir, inputName+ext)

	err = ioutil.WriteFile(inputPath, input, os.ModePerm)
	if err != nil {
		return "", err
	}

	outputBase := inputName + ".json"
	outputPath := filepath.Join(client.config.OutputDir, outputBase)
	return outputPath, nil
}

func (client *client) collectOutput(ctx context.Context, path string) (*DupeDetection, error) {
	// Waiting for JSON output from dupe detection service
	ctx, cancel := context.WithTimeout(ctx, time.Duration(client.config.WaitForOutputTimeout)*time.Second)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("timeout: %w", ctx.Err())
		case <-ticker.C:
			result, err := parseOutput(path)

			if err != nil {
				log.WithContext(ctx).WithError(err).Debug("Failed to collect output from dupe detection service")
				continue
			}
			return result, nil
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

	result := DupeDetection{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Errorf("failed to parse dupe detection service JSON output: %v", err)
	}

	var fingerprints []float64

	// Check fingerprints json string is right format
	err = json.Unmarshal([]byte(result.FingerPrints), &fingerprints)
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
