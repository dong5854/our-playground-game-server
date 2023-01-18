package server

import "net"

type TCPServer interface {
	net.Listener
	Run()
}
