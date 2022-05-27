package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/auth"
	"github.com/pastelnetwork/gonode/p2p/test/auth/common"
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

	var conn net.Conn
	var err error

	addr := fmt.Sprintf("%s:%s", ip, *port)

	// new the network socket
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("new socket: %v", err))
	}

	secInfo := &alts.SecInfo{
		PastelID: common.ServerB64PubKey,
	}
	fakePastelClient := common.NewSecClient()

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
			authenticator := kademlia.NewAuthHelper(fakePastelClient, secInfo)
			authHandshaker, _ := auth.NewServerHandshaker(context.Background(), authenticator, rawConn)
			conn, err = authHandshaker.ServerHandshake(context.Background())

			if err != nil {
				panic(fmt.Sprintf("client auth establish %q: %v", rawConn.RemoteAddr(), err))
			}

			for {
				data := make([]byte, 1024)
				rlen, err := conn.Read(data)
				if err != nil {
					fmt.Printf("socket read error: %v\n", err)
					return
				}

				fmt.Printf("received :%s\n", string(data[0:rlen]))
				if _, err := conn.Write(data[:rlen]); err != nil {
					fmt.Printf("socket write error: %v\n", err)
				}

			}
		}(rawConn)
	}
}
