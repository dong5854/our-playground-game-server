package server

import (
	"io"
	"log"
	"net"

	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/util/threadsafe"
)

const (
	TCP     = "tcp"
	MaxUser = 1000
)

type tcpServer struct {
	net.Listener
	handler.TCPHandler
	*threadsafe.TCPChannels
	*threadsafe.ClientList
}

func NewTCPServer(address string) TCPServer {
	server := new(tcpServer)
	var err error
	server.Listener, err = net.Listen(TCP, address)
	if err != nil {
		log.Panic(err)
	}
	server.TCPChannels = &threadsafe.TCPChannels{
		FromClient: make(chan []byte, MaxUser),
		ToClient:   make(chan []byte, MaxUser),
		ErrChan:    make(chan error, 1),
	}
	server.ClientList = new(threadsafe.ClientList)
	server.TCPHandler = handler.NewTCPHandler(server.TCPChannels, server.ClientList)
	return server
}

func (t *tcpServer) Run() {
	defer log.Println("Stopped TCPServer")
	log.Println("Start TCPServer")
	go t.HandlePacket() // 패킷 핸들링
	for {
		log.Println("waiting for TCP HandShake")
		conn, err := t.Accept()
		log.Println("successfully HandShaken")
		if err != nil {
			log.Panic(err)
		}
		t.ClientList.Append(conn)
		go t.ReadPacket(conn)
		if len(t.ErrChan) != 0 {
			log.Panic(<-t.ErrChan)
		}
	}
}

func (t *tcpServer) ReadPacket(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			t.ErrChan <- err
		}
	}
	t.FromClient <- buf[:n]
}
