package net

import (
	"bufio"
	"net"

	"github.com/cloudflare/circl/dh/x448"
	"github.com/pastelnetwork/gonode/common/errors"
)

const okResponse = "ok"

type handShake struct {
	Conn 				net.Conn
	HandshakeService 	Handshaker
}

type HandShake interface {
	ServerHandshake() (Cipher, error)
	ClientHandshake() (Cipher, error)
}


func (service *handShake) ServerHandshake() (Cipher, error) {
	reader := bufio.NewReader(service.Conn)
	_, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	ok, err := service.HandshakeService.ServerHello()
	if err != nil {
		return nil, err
	}
	if _, err := service.Conn.Write(prepareLine(ok)); err != nil {
		return nil, err
	}
	clientPublicKey := make([]byte, x448.Size)
	_, err = reader.Read(clientPublicKey)
	if err != nil {
		return nil, err
	}
	cipher, err := service.HandshakeService.ServerKeyVerify(clientPublicKey)
	if err != nil {
		return nil, err
	}
	serverKey, err := service.HandshakeService.ServerKeyExchange()
	if err != nil {
		return nil, err
	}
	if _, err := service.Conn.Write(prepareLine(serverKey)); err != nil {
		return nil, err
	}
	return cipher, nil
}

func (service *handShake) ClientHandshake() (Cipher, error) {
	cl, err := service.HandshakeService.ClientHello()
	if err != nil {
		return nil, err
	}
	_, err = service.Conn.Write(prepareLine(cl))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(service.Conn)
	serverOK, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if string(serverOK) != okResponse {
		return nil, errors.Errorf("not expected server response")
	}
	clientsKey, err := service.HandshakeService.ClientKeyExchange()
	if err != nil {
		return nil, err
	}
	_, err = service.Conn.Write(clientsKey)
	if err != nil {
		return nil, err
	}

	serverPublicKey := make([]byte, x448.Size)
	_, err = reader.Read(serverPublicKey)
	if err != nil {
		return nil, err
	}
	cipher, err := service.HandshakeService.ClientKeyVerify(serverPublicKey)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}


func NewHandShake(conn net.Conn, h Handshaker) HandShake{
	return &handShake {
		Conn: conn,
		HandshakeService: h,
	}
}
