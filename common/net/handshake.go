package net

import (
	"bufio"
	"fmt"
	"net"

	"github.com/cloudflare/circl/dh/x448"
)

func ServerHandshake(c net.Conn, h Handshaker) (Cipher, error) {
	reader := bufio.NewReader(c)
	_, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	ok, err := h.ServerHello()
	if _, err := c.Write(prepareLine(ok)); err != nil {
		return nil, err
	}
	clientPublicKey := make([]byte, x448.Size)
	_, err = reader.Read(clientPublicKey)
	if err != nil {
		return nil, err
	}
	cipher, err := h.ServerKeyVerify(clientPublicKey)
	if err != nil {
		return nil, err
	}
	serverKey, err := h.ServerKeyExchange()
	if _, err := c.Write(prepareLine(serverKey)); err != nil {
		return nil, err
	}
	return cipher, nil
}

func ClientHandshake(c net.Conn, h Handshaker) (Cipher, error) {
	cl, err := h.ClientHello()
	if err != nil {
		return nil, err
	}
	_, err = c.Write(prepareLine(cl))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(c)
	serverOK, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if string(serverOK) != "ok" {
		return nil, fmt.Errorf("not expected server response")
	}
	clientsKey, err := h.ClientKeyExchange()
	if err != nil {
		return nil, err
	}
	_, err = c.Write(clientsKey)
	if err != nil {
		return nil, err
	}

	serverPublicKey := make([]byte, x448.Size)
	_, err = reader.Read(serverPublicKey)
	if err != nil {
		return nil, err
	}
	cipher, err := h.ClientKeyVerify(serverPublicKey)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func prepareLine(r []byte) []byte {
	return append(r, []byte("\n")...)
}
