package dupedetection

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		config *Config
	}{
		{
			config: &Config{
				InputDir:  DefaultInputDir,
				OutputDir: DefaultOutputDir,
			},
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			config := NewConfig()
			assert.Equal(t, testCase.config, config)
		})
	}
}

func TestSetWorkDir(t *testing.T) {
	t.Parallel()

	type args struct {
		config  *Config
		workDir string
	}

	workDir, err := os.Getwd()
	assert.Equal(t, nil, err)

	absoluteInputDir := path.Join(workDir, DefaultInputDir)
	absoluteOutputDir := path.Join(workDir, DefaultOutputDir)

	testCases := []struct {
		args args
		want *Config
	}{
		{
			args: args{
				config: &Config{
					InputDir:  DefaultInputDir,
					OutputDir: DefaultOutputDir,
				},
				workDir: workDir,
			},
			want: &Config{
				InputDir:  absoluteInputDir,
				OutputDir: absoluteOutputDir,
			},
		},
		{
			args: args{
				config: &Config{
					InputDir:  absoluteInputDir,
					OutputDir: absoluteOutputDir,
				},
				workDir: workDir,
			},
			want: &Config{
				InputDir:  absoluteInputDir,
				OutputDir: absoluteOutputDir,
			},
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			config := testCase.args.config
			config.SetWorkDir(testCase.args.workDir)
			assert.Equal(t, testCase.want, config)
		})
	}
}
