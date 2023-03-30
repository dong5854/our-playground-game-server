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
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/packets"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

const (
	ChatTypeCHAT    = "CHAT"
	ChatTypeWHISPER = "WHISPER"
	ChatTypeSYSTEM  = "SYSTEM"
	MaxUser         = 1000
)

type tcpHandler struct {
	chatParser  packets.ChatParser
	clientMap   *sync.Map
	tcpChannels *threadsafe.TCPChannels
}

func NewTCPHandler(chatParser packets.ChatParser, tcpChannels *threadsafe.TCPChannels, clientMap *sync.Map) TCPHandler {
	return &tcpHandler{
		chatParser:  chatParser,
		tcpChannels: tcpChannels,
		clientMap:   clientMap,
	}
}

func (t *tcpHandler) TCPChannel() *threadsafe.TCPChannels {
	return t.tcpChannels
}

func (t *tcpHandler) HandlePacket() { // handlePacket 함수는 하나의 고루틴에서만 돌아감
	go t.readPacket() // 패킷을 읽어들이는 고루틴 하나 생성

	for { // 데이터를 받아와 채팅의 종류에 따라 처리하는 고루틴
		data := <-t.tcpChannels.FromClient
		logger.Debug("byte data: " + fmt.Sprint(data))

		t.chatParser.Parse(data)
		logger.Debug("Message: " + t.chatParser.Message())

		switch t.chatParser.Type() {
		case ChatTypeCHAT:
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
