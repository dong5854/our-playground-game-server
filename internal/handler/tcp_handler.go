package handler

import (
	"io"
	"net"
	"strconv"
	"sync"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/logger"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/parser"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

const (
	ECHO    = "echo"
	MaxUser = 1000
)

type tcpHandler struct {
	parser      parser.Parser
	clientMap   *sync.Map
	tcpChannels *threadsafe.TCPChannels
}

func NewTCPHandler(parser parser.Parser, tcpChannels *threadsafe.TCPChannels, ClientMap *sync.Map) TCPHandler {
	return &tcpHandler{
		parser:      parser,
		tcpChannels: tcpChannels,
		clientMap:   ClientMap,
	}
}

func (t *tcpHandler) TCPChannel() *threadsafe.TCPChannels {
	return t.tcpChannels
}

func (t *tcpHandler) HandlePacket() { // handlePacket 함수는 하나의 고루틴에서만 돌아감
	go t.readPacket() // 패킷을 읽어들이는 고루틴 하나 생성

	for { // 데이터를 받아와 데이터의 종류마다 다른 메소드로 핸들링.
		data := <-t.tcpChannels.FromClient
		logger.Debug("data: " + string(data))
		if err := t.parser.Unmarshal(data); err != nil {
			logger.Error(err.Error())
			t.tcpChannels.ErrChan <- err
		}
		logger.Debug("query: " + t.parser.Query())
		if t.parser.Query() == ECHO {
			go t.echoToAllClients(data)
		}
	}
}

func (t *tcpHandler) readPacket() {
	for { // 계속 실행되어야 하므로 무한 loop
		t.clientMap.Range(func(key, value any) bool {
			if conn, ok := value.(net.Conn); ok {
				buf := make([]byte, 1024)
				logger.Debug("waiting to read from id: " + key.(string))
				n, err := conn.Read(buf) // non-blocking
				logger.Debug("message read: length " + strconv.Itoa(n))

				if err != nil {
					if err != io.EOF {
						logger.Error("error on reading from connection from: " + key.(string))
						t.removeClient(key.(string), conn)
					}
				}

				if n > 0 { // 읽어들인 값이 없으면 채널에 값을 보내지 않음
					t.tcpChannels.FromClient <- buf[:n]
				}
			}
			return true
		})
	}
}

func (t *tcpHandler) echoToAllClients(data []byte) {
	t.clientMap.Range(func(key, value any) bool {
		if conn, ok := value.(net.Conn); ok {
			if _, err := conn.Write(data); err != nil {
				logger.Error("error on writing to connection to: " + key.(string))
				t.removeClient(key.(string), conn)
			}
		}
		return true
	})
}

func (t *tcpHandler) removeClient(uuid string, client net.Conn) {
	defer client.Close()
	t.clientMap.Delete(uuid)
}
