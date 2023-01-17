package handler

import (
	"log"
	"net"
	"sync"

	"github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/util/threadsafe"
)

type tcpHandler struct {
	clientMap *sync.Map
	*threadsafe.TCPChannels
}

func NewTCPHandler(channelSet *threadsafe.TCPChannels, ClientMap *sync.Map) TCPHandler {
	return &tcpHandler{
		TCPChannels: channelSet,
		clientMap:   ClientMap,
	}
}

func (t *tcpHandler) HandlePacket() {
	data := <-t.FromClient
	// if : data 가 echo 라면
	t.echoToAllClients(data)
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
