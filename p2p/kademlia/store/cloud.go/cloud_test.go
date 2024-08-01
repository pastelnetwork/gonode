package cloud

/*
import (
	"fmt"
	"os"
	"testing"
)

func TestRcloneStorage_Fetch(t *testing.T) {
	SetRcloneB2Env("0058a83f04e57580000000002", "K005lP8Wd0Gvr4JOdEI6e6BTAyt6iZA")
	storage := NewRcloneStorage("pastel-test", "b2")

	tests := map[string]struct {
		key      string
		expected string
		wantErr  bool
	}{
		"File Exists": {
			key:      "copy_45707.txt",
			expected: "expected content of testfile.txt",
			wantErr:  false,
		},
		"File Does Not Exist": {
			key:     "copy_45707_nonexistent.txt",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := storage.Fetch(tc.key)
			if (err != nil) != tc.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && string(data) != tc.expected {
				t.Errorf("Fetch() got = %v, want %v", string(data), tc.expected)
			}
		})
	}
}

// SetRcloneB2Env sets the environment variables for rclone configuration.
func SetRcloneB2Env(accountID, appKey string) error {
	err := os.Setenv("RCLONE_CONFIG_B2_TYPE", "b2")
	if err != nil {
		return fmt.Errorf("failed to set RCLONE_CONFIG_B2_TYPE: %w", err)
	}

	err = os.Setenv("RCLONE_CONFIG_B2_ACCOUNT", accountID)
	if err != nil {
		return fmt.Errorf("failed to set RCLONE_CONFIG_B2_ACCOUNT: %w", err)
	}

	err = os.Setenv("RCLONE_CONFIG_B2_KEY", appKey)
	if err != nil {
		return fmt.Errorf("failed to set RCLONE_CONFIG_B2_KEY: %w", err)
	}

	return nil
}


func TestRcloneStorage_FetchBatch(t *testing.T) {
	storage := NewRcloneStorage("mybucket")

	tests := map[string]struct {
		keys     []string
		expected map[string]string
		wantErr  bool
	}{
		"Multiple Files Exist": {
			keys: []string{"testfile1.txt", "testfile2.txt"}, // Ensure these files exist
			expected: map[string]string{
				"testfile1.txt": "content of testfile1.txt",
				"testfile2.txt": "content of testfile2.txt",
			},
			wantErr: false,
		},
		"Some Files Do Not Exist": {
			keys: []string{"testfile1.txt", "doesNotExist.txt"},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			results, err := storage.FetchBatch(tc.keys)
			if (err != nil) != tc.wantErr {
				t.Errorf("FetchBatch() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr {
				for k, v := range tc.expected {
					if results[k] != v {
						t.Errorf("FetchBatch() got = %v, want %v for key %v", string(results[k]), v, k)
					}
				}
			}
		})
	}
}
*/
