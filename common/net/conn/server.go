package conn

import (
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"net"
)

type Server struct {
	net.Listener
	transport transport.Transport
}

func NewServer(listener net.Listener, transport transport.Transport) *Server {
	return &Server{
		Listener:  listener,
		transport: transport,
	}
}

func (s *Server) Accept() (net.Conn, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if s.transport != nil {
		return s.transport.ServerHandshake(conn)
	}

	return conn, nil
}
