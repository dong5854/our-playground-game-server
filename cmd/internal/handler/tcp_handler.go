package handler

import (
	"io"
	"log"
	"net"
	"sync"

	idl "github.com/Team-OurPlayground/idl/proto"
	"google.golang.org/protobuf/proto"

	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/util/threadsafe"
)

const ECHO = "echo"

type tcpHandler struct {
	clientMap   *sync.Map
	tcpChannels *threadsafe.TCPChannels
}

func NewTCPHandler(tcpChannels *threadsafe.TCPChannels, ClientMap *sync.Map) TCPHandler {
	return &tcpHandler{
		tcpChannels: tcpChannels,
		clientMap:   ClientMap,
	}
}

func (t *tcpHandler) TCPChannel() *threadsafe.TCPChannels {
	return t.tcpChannels
}

func (t *tcpHandler) HandlePacket() {
	go t.readPacket() // 패킷을 읽어들이는 고루틴 하나 생성

	for { // 데이터를 받아와 데이터의 종류마다 다른 메소드로 핸들링.
		data := <-t.tcpChannels.FromClient

		message := &idl.SearchRequest{}
		if err := proto.Unmarshal(data, message); err != nil {
			t.tcpChannels.ErrChan <- err
		}

		if message.Query == ECHO {
			go t.echoToAllClients(data)
		}
	}
}

func (t *tcpHandler) readPacket() {
	for { // 계속 실행되어야 하므로 무한 loop
		t.clientMap.Range(func(key, value any) bool {
			if conn, ok := value.(net.Conn); ok {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf) // non-blocking

				if err != nil {
					if err != io.EOF {
						log.Println("error on reading from connection")
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
				log.Println("error on writing to connection")
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
