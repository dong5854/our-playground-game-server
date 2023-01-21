package main

import (
	"sync"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/server"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/parser"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

func main() {

	parser := parser.NewProtobufParser()
	tcpChannels := &threadsafe.TCPChannels{
		FromClient: make(chan []byte, handler.MaxUser),
		ToClient:   make(chan []byte, handler.MaxUser),
		ErrChan:    make(chan error, 1),
	}
	tcpHandler := handler.NewTCPHandler(parser, tcpChannels, new(sync.Map))
	server := server.NewTCPServer("127.0.0.1:6112", tcpHandler)
	server.Run()
	defer server.Close()
}
