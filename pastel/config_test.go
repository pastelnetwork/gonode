package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPort(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg          *Config
		wantPort     int
		wantUsername string
		wantPassword string
		wantHostname string
	}{
		"simple": {
			cfg:          NewConfig(),
			wantPort:     defaultMainnetPort,
			wantHostname: defaultHostname,
		},
		"testnet-on": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:     0,
					Testnet:  1,
					Username: "rpc",
					Hostname: "test",
				},
			},
			wantUsername: "rpc",
			wantHostname: "test",
			wantPort:     defaultTestnetPort,
		},
		"test-username": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:     0,
					Testnet:  1,
					Username: "rpc",
				},
			},
			wantUsername: "rpc",
			wantPort:     defaultTestnetPort,
		},
		"test-hostname": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:     0,
					Testnet:  1,
					Hostname: "test",
				},
			},
			wantHostname: "test",
			wantPort:     defaultTestnetPort,
		},
		"test-password": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:     0,
					Testnet:  1,
					Password: "test",
				},
			},
			wantPassword: "test",
			wantPort:     defaultTestnetPort,
		},
		"testnet-off": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:    0,
					Testnet: 0,
				},
			},
			wantPort: defaultMainnetPort,
		},
		"testnet-on-port-provided": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:    12170,
					Testnet: 1,
				},
			},
			wantPort: 12170,
		},
		"testnet-off-port-provided": {
			cfg: &Config{
				ExternalConfig: &ExternalConfig{
					Port:    12170,
					Testnet: 0,
				},
			},
			wantPort: 12170,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.cfg.port()
			assert.Equal(t, tc.wantPort, got)

			gotVal := tc.cfg.username()
			assert.Equal(t, tc.wantUsername, gotVal)

			gotVal = tc.cfg.hostname()
			assert.Equal(t, tc.wantHostname, gotVal)

			gotVal = tc.cfg.password()
			assert.Equal(t, tc.wantPassword, gotVal)
		})
	}
}
