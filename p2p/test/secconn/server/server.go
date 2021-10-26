package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/anacrolix/utp"
	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/test/secconn/common"
)

func main() {
	var (
		ip   = "0.0.0.0"
		port = flag.String("port", "", "port to connect")
	)
	flag.Parse()

	if *port == "" {
		flag.PrintDefaults()
		return
	}

	priKey, pubKey := common.GetKeys()

	var conn net.Conn
	var err error

	addr := fmt.Sprintf("%s:%s", ip, *port)

	// new the network socket
	socket, err := utp.NewSocket("udp", addr)
	if err != nil {
		panic(fmt.Sprintf("new socket: %v", err))
	}

	secInfo := &alts.SecInfo{}
	fakePastelClient := &common.SecClient{
		Curve: ed448.NewCurve(),
		Pri:   priKey,
		Pub:   pubKey,
	}
	transportCredentials := credentials.NewClientCreds(fakePastelClient, secInfo)

	for {

		// accept the incomming connections
		rawConn, err := socket.Accept()
		if err != nil {
			fmt.Printf("socket accept: %v\n", err)
			continue
		}

		go func(rawConn net.Conn) {
			defer rawConn.Close()
			// process request
			conn, _, err = transportCredentials.ServerHandshake(rawConn)
			if err != nil {
				fmt.Printf("socket handshaking: %v\n", err)
				return
			}

			for {
				data := make([]byte, 1024)
				len, err := conn.Read(data)
				if err != nil {
					fmt.Printf("socket read: %v\n", err)
					return
				}

				fmt.Printf("received :%s\n", string(data[0:len]))
			}
		}(rawConn)
	}
}
