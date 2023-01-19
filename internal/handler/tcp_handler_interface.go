package handler

import (
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

type TCPHandler interface {
	HandlePacket()
	TCPChannel() *threadsafe.TCPChannels
}
