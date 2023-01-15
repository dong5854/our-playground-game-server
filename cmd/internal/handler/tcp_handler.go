package handler

import "github.com/Team-OurPlayground/our-playground-game-server/cmd/internal/util/threadsafe"

type tcpHandler struct {
	*threadsafe.TCPChannels
	*threadsafe.ClientList
}

func NewTCPHandler(channelSet *threadsafe.TCPChannels, ClientList *threadsafe.ClientList) TCPHandler {
	return &tcpHandler{
		TCPChannels: channelSet,
		ClientList:  ClientList,
	}
}

func (t *tcpHandler) HandlePacket() {
	data := <-t.FromClient
	// if : data 가 echo 라면
	t.echoToAllClients(data)
}

func (t *tcpHandler) echoToAllClients(data []byte) {
	t.ClientList.RLock()
	defer t.ClientList.RUnlock()
	clients := t.ClientList.Get()
	for _, client := range clients {
		_, err := client.Write(data)
		if err != nil {
			t.ErrChan <- err
			break
		}
	}
}
