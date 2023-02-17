package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Team-OurPlayground/idl/goproto"
	"google.golang.org/protobuf/proto"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/structs"
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

		go t.authenticatePlayer(conn)

		if len(t.ErrChan) != 0 {
			err = <-t.ErrChan
			logger.Error(err.Error())
			log.Panic(err)
		}
	}
}

func (t *tcpServer) authenticatePlayer(conn net.Conn) {
	logger.Debug("start player authenticate process")
	tokenByte := make([]byte, 1024) // 인증 토큰 값 읽어옴
	n, err := conn.Read(tokenByte)
	if err != nil {
		if err != io.EOF {
			logger.Error("error on reading tokenByte from connection")
			conn.Close()
			return
		}
	}
	logger.Debug("tokenByte read: length" + strconv.Itoa(n))
	AuthInfo := new(goproto.Authenticate)
	if err := proto.Unmarshal(tokenByte[:n], AuthInfo); err != nil {
		logger.Error("error unMarshaling authentication info")
		conn.Close()
		return
	}
	// TODO: jwt 토큰으로 인증하는 프로세스 개발 필요, 클라이언트 개발 편의성을 위해 후순위
	logger.Debug("authentication success")
	id := AuthInfo.Id
	player := structs.Player{
		ID:   AuthInfo.Id,
		Conn: conn,
	}
	t.clientMap.Store(id, player)
	logger.Debug("client connected as id: " + id)
}
