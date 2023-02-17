package server

import "net"

type Room struct {
	id     string
	player []*Player
}

type Player struct {
	id   string
	conn net.Conn
}
