package grpc

type Client interface {
	Connect(host string)
}
