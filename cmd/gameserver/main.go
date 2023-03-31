package main

import (
	"sync"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/server"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/logger"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/packets"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

func main() {
	// TODO: heartbeat 추가하기
	chatParser := packets.NewChatParser()
	tcpChannels := &threadsafe.TCPChannels{
		FromClient: make(chan []byte, handler.MaxUser),
		ToClient:   make(chan []byte, handler.MaxUser),
		ErrChan:    make(chan error, 1),
	}
	clientMap := new(sync.Map)

	tcpHandler := handler.NewTCPHandler(chatParser, tcpChannels, clientMap)
	server := server.NewTCPServer("0.0.0.0:6112", tcpHandler, clientMap)
	server.Run()

	defer func() {
		server.Close()
		logger.Sync()
	}()
}
