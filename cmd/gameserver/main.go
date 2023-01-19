package main

import (
	"sync"

	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/server"
	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/util/threadsafe"
)

func main() {
	tcpHandler := handler.NewTCPHandler(new(threadsafe.TCPChannels), new(sync.Map))
	server := server.NewTCPServer("127.0.0.1:6112", tcpHandler)
	server.Run()
	defer server.Close()
}
