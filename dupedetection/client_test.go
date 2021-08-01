package dupedetection

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	type args struct {
		config *Config
	}
	testCases := []struct {
		args args
		want Client
	}{
		{
			args: args{
				config: &Config{
					InputDir:  "",
					OutputDir: "",
				},
			},
			want: &client{
				config: &Config{
					InputDir:  "",
					OutputDir: "",
				},
			},
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			client := NewClient(testCase.args.config)
			assert.Equal(t, testCase.want, client)
		})
	}
}

func TestCopyImageToInputDir(t *testing.T) {
	t.Parallel()

	type args struct {
		config   *Config
		testFile string
	}

	pwd, err := os.Getwd()
	assert.Equal(t, nil, err)
	inputDir := path.Join(pwd, "input")
	outputDir := path.Join(pwd, "output")
	err = os.Mkdir(inputDir, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.RemoveAll(inputDir)

	testFile := path.Join(pwd, "test.jpg")
	testData := []byte("test")
	err = ioutil.WriteFile(testFile, testData, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.Remove(testFile)

	testCases := []struct {
		args             args
		outputBaseExt    string
		outputPathPrefix string
		assertion        assert.ErrorAssertionFunc
	}{
		{
			args: args{
				testFile: testFile,
				config: &Config{
					InputDir:  inputDir,
					OutputDir: outputDir,
				},
			},
			outputBaseExt:    ".json",
			outputPathPrefix: outputDir,
			assertion:        assert.NoError,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			client := NewClient(testCase.args.config).(*client)

			outputPath, err := client.copyImageToInputDir(testCase.args.testFile)
			testCase.assertion(t, err)
			assert.True(t, filepath.Ext(outputPath) == testCase.outputBaseExt)
			assert.True(t, strings.HasPrefix(outputPath, testCase.outputPathPrefix))
		})
	}
}

func TestCollectOutput(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		config     *Config
		outputPath string
	}

	pwd, err := os.Getwd()
	assert.Equal(t, nil, err)

	inputDir := path.Join(pwd, "input")
	outputDir := path.Join(pwd, "output")
	err = os.Mkdir(outputDir, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.RemoveAll(outputDir)

	baseName := "test.json"
	outputPath := path.Join(outputDir, baseName)
	baseName1 := "test1.json"
	outputPath1 := path.Join(outputDir, baseName1)
	baseName2 := "test2.json"
	outputPath2 := path.Join(outputDir, baseName2)
	baseName3 := "test3.json"
	outputPath3 := path.Join(outputDir, baseName3)

	result := DupeDetection{
		DupeDetectionSystemVer:  "1.0",
		ImageHash:               "02bcagagaeagdae",
		PastelRarenessScore:     0.9,
		InternetRarenessScore:   0.8,
		MatchesFoundOnFirstPage: 1,
		NumberOfResultPages:     10,
		FirstMatchURL:           "https://example.com",
		OpenNSFWScore:           0.6,
		AlternateNSFWScores: AlternateNSFWScores{
			Drawings: 0.1,
			Hentai:   0.2,
			Neutral:  0.3,
			Porn:     0.4,
			Sexy:     0.5,
		},
		ImageHashes: ImageHashes{
			PerceptualHash: "c999d3d3230724fc",
			AverageHash:    "ffff990999181800",
			DifferenceHash: "7333237333337331",
		},
		FingerPrints: "[0.1, 0.2, 0.3]",
	}

	data, err := json.Marshal(&result)
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile(outputPath, data, os.ModePerm)
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile(outputPath2, []byte("test"), os.ModePerm)
	assert.Equal(t, nil, err)

	result1 := result
	result1.FingerPrints = "[0.1, 0.2, 0.3"
	data, err = json.Marshal(&result1)
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile(outputPath3, data, os.ModePerm)
	assert.Equal(t, nil, err)

	testCases := []struct {
		args      args
		result    *DupeDetection
		assertion assert.ErrorAssertionFunc
	}{
		{
			args: args{
				ctx: context.Background(),
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
				outputPath: outputPath,
			},
			result:    &result,
			assertion: assert.NoError,
		},
		{
			args: args{
				ctx: context.Background(),
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
				outputPath: outputPath1,
			},
			result:    nil,
			assertion: assert.Error,
		},
		{
			args: args{
				ctx: context.Background(),
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
				outputPath: outputPath2,
			},
			result:    nil,
			assertion: assert.Error,
		},
		{
			args: args{
				ctx: context.Background(),
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
				outputPath: outputPath3,
			},
			result:    nil,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			client := NewClient(testCase.args.config).(*client)

			ctx := context.Background()

			result, err := client.collectOutput(ctx, testCase.args.outputPath)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.result, result)
		})
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		config   *Config
		testFile string
	}

	pwd, err := os.Getwd()
	assert.Equal(t, nil, err)
	inputDir := path.Join(pwd, "geninput")
	outputDir := path.Join(pwd, "genoutput")
	err = os.Mkdir(inputDir, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.RemoveAll(inputDir)
	err = os.Mkdir(outputDir, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.RemoveAll(outputDir)

	testFile := path.Join(pwd, "test.jpg")
	testFile2 := path.Join(pwd, "test2.jpg")
	testData := []byte("test")
	err = ioutil.WriteFile(testFile, testData, os.ModePerm)
	assert.Equal(t, nil, err)
	defer os.Remove(testFile)

	testCases := []struct {
		args      args
		result    *DupeDetection
		assertion assert.ErrorAssertionFunc
	}{
		{
			args: args{
				ctx:      context.Background(),
				testFile: testFile,
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
			},
			result:    nil,
			assertion: assert.Error,
		},
		{
			args: args{
				ctx:      context.Background(),
				testFile: testFile2,
				config: &Config{
					InputDir:             inputDir,
					OutputDir:            outputDir,
					WaitForOutputTimeout: 1,
				},
			},
			result:    nil,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			client := NewClient(testCase.args.config).(*client)

			ctx := context.Background()

			result, err := client.Generate(ctx, testCase.args.testFile)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.result, result)
		})
	}
}
