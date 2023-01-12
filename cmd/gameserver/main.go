package main

import (
	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/server"
)

func main() {
	server := server.NewTCPServer("127.0.0.1:6112")
	server.Run()
	defer server.Close()
}
