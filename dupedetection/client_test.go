package dupedetection

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

			outputBase, outputPath, err := client.copyImageToInputDir(testCase.args.testFile)
			testCase.assertion(t, err)
			assert.True(t, filepath.Ext(outputBase) == testCase.outputBaseExt)
			assert.True(t, strings.HasPrefix(outputPath, testCase.outputPathPrefix))
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
				ctx:      context.Background(),
				testFile: testFile,
				config: &Config{
					InputDir:  inputDir,
					OutputDir: outputDir,
				},
			},
			outputBaseExt:    ".json",
			outputPathPrefix: outputDir,
			assertion:        assert.Error,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			client := NewClient(testCase.args.config).(*client)

			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()

			result, err := client.Generate(ctx, testCase.args.testFile)
			testCase.assertion(t, err)
			assert.Equal(t, (*DupeDetectionResult)(nil), result)
		})
	}
}
