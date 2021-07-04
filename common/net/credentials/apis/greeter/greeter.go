package greeter

import context "context"

type Greeter struct {
	UnimplementedGreeterServiceServer
}

func (s *Greeter) SayHi(_ context.Context, request *SayHiRequest) (*SayHiReply, error) {
	return &SayHiReply{Message: request.Name}, nil
}
