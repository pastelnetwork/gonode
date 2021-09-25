package main

func main() {

}

/*
import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/dupedetection"
)

const (
	defaultTestImage = "test.png"
)

func main() {
	fmt.Println("Start dupe detection test...")
	defer fmt.Println("Finished dupe detection test.")

	workDir := flag.String("workDir", "", "a path to the working directory.")
	flag.Parse()

	currentDir, _ := os.Getwd()
	if *workDir != "" {
		currentDir = *workDir
	}

	ddconf := dupedetection.NewConfig()
	ddconf.SetWorkDir(currentDir)

	fmt.Printf("input directory: %s\n", ddconf.InputDir)
	fmt.Printf("output directory: %s\n", ddconf.OutputDir)
	client := dupedetection.NewClient(ddconf)

	// read image content
	img, err := ioutil.ReadFile(defaultTestImage)
	if err != nil {
		fmt.Printf("Could not read image content: %v\n", err)
		return
	}

	format := filepath.Ext(defaultTestImage)

	ctx := context.Background()
	result, err := client.Generate(ctx, img, strings.TrimPrefix(format, "."))
	if err != nil {
		fmt.Printf("Could not get fringerprints from dupe detection service: %v\n", err)
		return
	}
	fmt.Printf("Dupe detection system version: %s\n", result.DupeDetectionSystemVer)
	fmt.Printf("Image hash: %s\n", result.ImageHash)
	fmt.Printf("Pastel Rareness Score: %v\n", result.PastelRarenessScore)
	fmt.Printf("Internet Rareness Score: %v\n", result.InternetRarenessScore)
	fmt.Printf("Open NSFW Score: %v\n", result.OpenNSFWScore)
	fmt.Printf("Matches found on first page: %v\n", result.MatchesFoundOnFirstPage)
	fmt.Printf("Number of result pages: %v\n", result.NumberOfResultPages)
	fmt.Printf("First match URL: %s\n", result.FirstMatchURL)
	fmt.Printf("Alternate NSFW Scores: %v\n", result.AlternateNSFWScores)
	fmt.Printf("Image Hashes: %v\n", result.ImageHashes)
	fmt.Printf("Fingerprints: %v\n", result.Fingerprints)
}
*/
