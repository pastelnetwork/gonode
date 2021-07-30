package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/dupedetection"
)

const (
	defaultConfigFile = "./test.yml"
	defaultTestImage  = "test.png"
)

type DDConfig struct {
	DupeDection *dupedetection.Config `mapstructure:"dd-service" json:"dd-service,omitempty"`
}

func main() {
	fmt.Println("Start dupe detection test...")
	defer fmt.Println("Finished dupe detection test.")

	workDir := flag.String("workDir", "", "a path to the working directory.")
	flag.Parse()

	currentDir, _ := os.Getwd()
	if *workDir != "" {
		currentDir = *workDir
	}
	config.DupeDection.SetWorkDir(currentDir)

	ddconf := dupedetection.NewConfig()
	ddconf.SetWorkDir(currentDir)

	fmt.Printf("input directory: %s\n", ddconf.InputDir)
	fmt.Printf("output directory: %s\n", ddconf.OutputDir)
	client := dupedetection.NewClient(ddconf)

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
