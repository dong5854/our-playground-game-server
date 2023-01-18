package server

import (
	"log"
	"net"
	"sync"

	"github.com/google/uuid"

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
	clientMap *sync.Map
	*threadsafe.TCPChannels
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
	server.clientMap = &sync.Map{}
	server.TCPHandler = handler.NewTCPHandler(server.TCPChannels, server.clientMap)
	return server
}

func (t *tcpServer) Run() {
	defer log.Println("Stopped TCPServer")
	log.Println("Start TCPServer")
	go t.HandlePacket() // 패킷 핸들링 고루틴은 하나만 생성
	for {
		log.Println("waiting for TCP HandShake")
		conn, err := t.Accept()
		log.Println("successfully HandShaken")
		if err != nil {
			log.Panic(err)
		}

		id := uuid.New().String()
		t.clientMap.Store(id, conn)

		if len(t.ErrChan) != 0 {
			log.Panic(<-t.ErrChan)
		}
	}
}
