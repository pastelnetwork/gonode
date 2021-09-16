package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPort(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg      *Config
		wantPort int
	}{
		"simple": {
			cfg:      NewConfig(),
			wantPort: defaultMainnetPort,
		},
		"testnet-on": {
			cfg: &Config{
				Port:    0,
				Testnet: 1,
			},
			wantPort: defaultTestnetPort,
		},
		"testnet-off": {
			cfg: &Config{
				Port:    0,
				Testnet: 0,
			},
			wantPort: defaultMainnetPort,
		},
		"testnet-on-port-provided": {
			cfg: &Config{

				Port:    12170,
				Testnet: 1,
			},
			wantPort: 12170,
		},
		"testnet-off-port-provided": {
			cfg: &Config{
				Port:    12170,
				Testnet: 0,
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
		})
	}
}
