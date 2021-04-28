package grpc

// Client represents grpc client
type Client interface {
	Connect(host string) error
}
