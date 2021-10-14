package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/anacrolix/utp"
	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/test/secconn/common"
)

const defaultConnDeadline = 5 * time.Second

func main() {
	var (
		ip   = flag.String("ip", "", "ip of server")
		port = flag.String("port", "", "port to connect")
	)
	flag.Parse()

	if len(*ip) == 0 || len(*port) == 0 {
		flag.PrintDefaults()
		return
	}

	priKey, pubKey := common.GetKeys()

	var conn net.Conn
	var err error

	remoteAddr := fmt.Sprintf("%s:%s", *ip, *port)
	ctx := context.Background()

	// dial the remote address with udp network
	rawConn, err := utp.DialContext(ctx, remoteAddr)
	if err != nil {
		panic(fmt.Sprintf("dial %q: %v", remoteAddr, err))
	}
	defer rawConn.Close()

	// set the deadline for read and write
	rawConn.SetDeadline(time.Now().Add(defaultConnDeadline))

	secInfo := &alts.SecInfo{}
	fakePastelClient := &common.SecClient{
		Curve: ed448.NewCurve(),
		Pri:   priKey,
		Pub:   pubKey,
	}
	transportCredentials := credentials.NewClientCreds(fakePastelClient, secInfo)

	conn, _, err = transportCredentials.ClientHandshake(ctx, "", rawConn)
	if err != nil {
		panic(fmt.Sprintf("client secure establish %q: %v", remoteAddr, err))
	}

	// encode and send the request message
	data := []byte("Hello From Client")
	if err != nil {
		panic(fmt.Sprintf("encode: %v", err))
	}
	if _, err := conn.Write(data); err != nil {
		panic(fmt.Sprintf("conn write: %v", err))
	}

	fmt.Println("Send Done")
}
