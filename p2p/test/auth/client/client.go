package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"github.com/pastelnetwork/gonode/p2p/test/auth/common"
)

const defaultConnDeadline = 5 * time.Second

func main() {
	var (
		ip   = flag.String("ip", "", "ip of server")
		port = flag.String("port", "", "port to connect")
	)
	flag.Parse()

	if *ip == "" || *port == "" {
		flag.PrintDefaults()
		return
	}

	var conn net.Conn
	var err error

	remoteAddr := fmt.Sprintf("%s:%s", *ip, *port)
	ctx := context.Background()

	// dial the remote address with udp network
	var d net.Dialer
	rawConn, err := d.DialContext(ctx, "tcp", remoteAddr)
	if err != nil {
		panic(fmt.Sprintf("dial %q: %v", remoteAddr, err))
	}
	defer rawConn.Close()

	// set the deadline for read and write
	rawConn.SetDeadline(time.Now().Add(defaultConnDeadline))

	secInfo := &alts.SecInfo{
		PastelID: common.ClientB64PubKey,
	}
	fakePastelClient := common.NewSecClient()

	authenticator := kademlia.NewAuthHelper(fakePastelClient, secInfo)
	authHandshaker, _ := auth.NewClientHandshaker(ctx, authenticator, rawConn)
	conn, err = authHandshaker.ClientHandshake(ctx)

	if err != nil {
		panic(fmt.Sprintf("client auth establish %q: %v", remoteAddr, err))
	}

	// encode and send the request message
	wdata := []byte("Hello From Client")
	if err != nil {
		panic(fmt.Sprintf("encode: %v", err))
	}
	if _, err := conn.Write(wdata); err != nil {
		panic(fmt.Sprintf("conn write: %v", err))
	}

	rdata := make([]byte, 1024)
	if _, err := conn.Read(rdata); err != nil {
		panic(fmt.Sprintf("conn read: %v", err))
	}
}
