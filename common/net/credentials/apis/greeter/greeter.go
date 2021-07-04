package greeter

import context "context"

// Greeter implentation
type Greeter struct {
	UnimplementedGreeterServiceServer
}

// SayHi implements the api of Greeter service
func (s *Greeter) SayHi(_ context.Context, request *SayHiRequest) (*SayHiReply, error) {
	return &SayHiReply{Message: request.Name}, nil
}
