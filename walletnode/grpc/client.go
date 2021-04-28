package grpc

// Client represents grpc client
type client struct {
	config *Config
}

func (client *client) Connect(host string) error {
	return nil
}

// NewClient returns a new Client instance.
func NewClient(config *Config) Client {
	return &client{
		config: config,
	}
}
