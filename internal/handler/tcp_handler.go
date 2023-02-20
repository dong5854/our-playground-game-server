package handler

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/structs"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/logger"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/parser"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

const (
	SetID        = "setID"
	SimulateMove = "simulateMove"
	SetPosition  = "setPosition"
	MaxUser      = 1000
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
		logger.Debug("byte data: " + fmt.Sprint(data))

		if err := t.parser.Unmarshal(data); err != nil {
			logger.Error(err.Error())
			t.tcpChannels.ErrChan <- err
		}
		logger.Debug("function: " + t.parser.Function())
		logger.Debug("data: " + t.parser.Data())

		switch t.parser.Function() {
		case SimulateMove:
			fallthrough
		case SetPosition:
			go t.echoToAllClients(data)
		default:
			logger.Error("undefined function")
		}
	}
}

func (t *tcpHandler) readPacket() {
	for { // 계속 실행되어야 하므로 무한 loop
		t.clientMap.Range(func(key, value any) bool {
			if player, ok := value.(structs.Player); ok {
				buf := make([]byte, 1024)
				logger.Debug("waiting to read from id: " + key.(string))
				err := player.Conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
				if err != nil {
					logger.Error("error while setting deadline to connection: " + key.(string))
					t.removeClient(key.(string), player)
					return true
				}
				n, err := player.Conn.Read(buf) // non-blocking, 여기서 멈춤 여기 context deadline 추가
				logger.Debug("message read: length " + strconv.Itoa(n))

				if err != nil {
					if os.IsTimeout(err) {
						return true
					}
					if err != io.EOF {
						logger.Error("error on reading from connection from: " + key.(string) + err.Error())
						t.removeClient(key.(string), player)
					}
				}

				if n > 0 { // 읽어들인 값이 없으면 채널에 값을 보내지 않음
					logger.Debug("send data to tcpChannels.FromClient")
					t.tcpChannels.FromClient <- buf[:n]
				}
			}
			return true
		})
	}
}

func (t *tcpHandler) echoToAllClients(data []byte) {
	t.clientMap.Range(func(key, value any) bool {
		if player, ok := value.(structs.Player); ok {
			if _, err := player.Conn.Write(data); err != nil {
				logger.Error("error on writing to connection to: " + key.(string))
				t.removeClient(key.(string), player)
			}
		}
		return true
	})
}

func (t *tcpHandler) removeClient(id string, client structs.Player) {
	defer client.Conn.Close()
	t.clientMap.Delete(id)
}
