package cleaner

import "github.com/pastelnetwork/gonode/hermes/service/node"

// Config is the config struct for cleaner-service
type Config struct {
	snHost string
	snPort int

	sn node.SNClientInterface
}

// NewConfig sets the configurations for cleaner service
func NewConfig(host string, port int, SN node.SNClientInterface) Config {
	return Config{
		snHost: host,
		snPort: port,
		sn:     SN,
	}
}
