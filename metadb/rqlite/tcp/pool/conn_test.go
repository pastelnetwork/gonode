package pool

import (
	"net"
	"testing"
)

func TestConn_Impl(_ *testing.T) {
	var _ net.Conn = new(Conn)
}
