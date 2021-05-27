package conn

import (
	"context"
	"github.com/gonode/common/net/conn/transport"
	"net"
)

type Server struct {
	net.Listener
	context   context.Context
	transport transport.Transport
}

func NewServer(listener net.Listener, context context.Context, transport transport.Transport) *Server {
	return &Server{
		Listener:  listener,
		transport: transport,
		context:   context,
	}
}

func (s *Server) Accept() (net.Conn, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if s.transport != nil {
		return s.transport.ServerHandshake(s.context, conn)
	}

	return conn, nil
}
