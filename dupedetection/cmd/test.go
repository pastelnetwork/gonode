package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/dupedetection"
)

const (
	defaultConfigFile = "./test.yml"
	defaultTestImage  = "./test.png"
)

type DDConfig struct {
	DupeDection *dupedetection.Config `mapstructure:"dd-service" json:"dd-service,omitempty"`
}

func main() {
	workDir := flag.String("workDir", "", "a path to the working directory.")
	confFile := flag.String("confFile", "", "a path to the configuration file.")
	flag.Parse()

	configPath := defaultConfigFile
	if *confFile != "" {
		configPath = *confFile
	}
	config := DDConfig{}
	if err := configurer.ParseFile(configPath, &config); err != nil {
		fmt.Errorf("could not parse config file: %v\n", err)
		return
	}

	currentDir, _ := os.Getwd()
	if *workDir != "" {
		currentDir = *workDir
	}
	config.DupeDection.SetWorkDir(currentDir)

	client := dupedetection.NewClient(config.DupeDection)

	ctx := context.Background()
	result, err := client.Generate(ctx, defaultTestImage)
	if err != nil {
		fmt.Errorf("could not get fringerprints from dupe detection service: %v\n", err)
		return
	}
	fmt.Printf("Pastel Rareness Score: %v\n", result.PastelRarenessScore)
	fmt.Printf("Internet Rareness Score: %v\n", result.InternetRarenessScore)
	fmt.Printf("Open NSFW Score: %v\n", result.OpenNSFWScore)
	fmt.Printf("Alternate NSFW Scores: %v\n", result.AlternateNSFWScores)
	fmt.Printf("Fingerprints: %v\n", result.FingerPrints)
}
