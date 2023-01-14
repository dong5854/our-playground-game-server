package server

import "net"

type TCPServer interface {
	net.Listener
	Run()
	ReadPacket(conn net.Conn)
}
