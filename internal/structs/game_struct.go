package structs

import "net"

type Room struct {
	ID     string
	Player []*Player
}

type Player struct {
	ID   string
	Conn net.Conn
}
