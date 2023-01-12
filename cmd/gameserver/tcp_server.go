package main

import (
	"io"
	"log"
	"net"
)

const TCP = "tcp"

type tcpServer struct {
	tcpListener net.Listener
}

func newTCPServer(address string) *tcpServer {
	server := new(tcpServer)
	var err error
	server.tcpListener, err = net.Listen(TCP, address)
	if err != nil {
		log.Panic(err)
	}
	return server
}

func (t *tcpServer) accept() (net.Conn, error) {
	conn, err := t.tcpListener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (t *tcpServer) readPacket(c net.Conn) ([]byte, error) {
	defer c.Close()
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	return buf[:n], err
}

func (t *tcpServer) close() {
	t.tcpListener.Close()
}
