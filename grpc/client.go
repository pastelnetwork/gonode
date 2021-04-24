package grpc

// Client represents grpc client
type client struct {
	config *Config
}

// NewClient returns a new Client instance.
func NewClient(config *Config) Client {
	return &client{
		config: config,
	}
}
