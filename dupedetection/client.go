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

// Generate implements dupedection.Client.Generate
func (client *client) Generate(ctx context.Context, path string) (*DupeDetectionResult, error) {

	// Copy iamge file to dupe detection service input directory
	outputBase, outputPath, err := client.copyImageToInputDir(path)
	if err != nil {
		return nil, errors.Errorf("failed to copy image to dupe detection input directory: %w", err)
	}

	result, err := client.collectOutput(ctx, outputBase, outputPath)

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

func (client *client) collectOutput(ctx context.Context, baseName, path string) (*DupeDetectionResult, error) {
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
					jsonFile, err := os.Open(path)
					if err != nil {
						err = errors.Errorf("failed to open JSON output: %v", err)
						continue
					}
					defer jsonFile.Close()

					byteValue, err := ioutil.ReadAll(jsonFile)
					if err != nil {
						err = errors.Errorf("failed to read JSON output: %v", err)
						continue
					}

					result := DupeDetectionResult{}
					err = json.Unmarshal(byteValue, &result)
					if err != nil {
						err = errors.Errorf("failed to parse JSON output: %v", err)
						continue
					}

					return &result, nil
				}
			}
		}
	}
}

// NewClient return a new Client instance
func NewClient(config *Config) Client {
	return &client{
		config: config,
	}
}
