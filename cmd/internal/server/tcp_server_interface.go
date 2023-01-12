package server

import "net"

type TCPServer interface {
	Run()
	Accept() (net.Conn, error)
	ReadPacket(conn net.Conn)
	Close() error
}
