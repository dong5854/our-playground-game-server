package server

import (
	"log"
	"net"
	"sync"

	"github.com/google/uuid"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/logger"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

const (
	TCP = "tcp"
)

type tcpServer struct {
	net.Listener
	handler.TCPHandler
	clientMap *sync.Map
	*threadsafe.TCPChannels
}

func NewTCPServer(address string, handler handler.TCPHandler, clientMap *sync.Map) TCPServer {
	server := new(tcpServer)
	var err error
	server.Listener, err = net.Listen(TCP, address)
	if err != nil {
		logger.Error(err.Error())
		log.Panic(err)
	}
	server.TCPChannels = handler.TCPChannel()
	server.clientMap = clientMap
	server.TCPHandler = handler
	return server
}

func (t *tcpServer) Run() {
	defer logger.Info("Stopped TCPServer")
	logger.Info("Start TCPServer")
	go t.HandlePacket() // 패킷 핸들링 고루틴은 하나만 생성
	for {
		logger.Info("waiting for TCP HandShake")
		conn, err := t.Accept()
		logger.Info("successfully HandShaken")
		if err != nil {
			logger.Error(err.Error())
			log.Panic(err)
		}

		id := uuid.New().String()
		t.clientMap.Store(id, conn)
		logger.Debug("client connected as id: " + id)

		if len(t.ErrChan) != 0 {
			err = <-t.ErrChan
			logger.Error(err.Error())
			log.Panic(err)
		}
	}
}
