package server_test

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Team-OurPlayground/idl/FBPacket"
	"github.com/Team-OurPlayground/idl/goproto"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/server"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/packets"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

type tcpServerSuite struct {
	suite.Suite
	server      server.TCPServer
	chatCreator packets.ChatCreator
	finishChan  chan struct{}
	once        sync.Once
}

func (suite *tcpServerSuite) SetupSuite() {
	suite.finishChan = make(chan struct{})
	suite.chatCreator = packets.NewChatCreator()
	chatParser := packets.NewChatParser()
	tcpChannels := &threadsafe.TCPChannels{
		FromClient: make(chan []byte, handler.MaxUser),
		ToClient:   make(chan []byte, handler.MaxUser),
		ErrChan:    make(chan error, 1),
	}
	clientMap := new(sync.Map)

	tcpHandler := handler.NewTCPHandler(chatParser, tcpChannels, clientMap)
	suite.server = server.NewTCPServer("127.0.0.1:64202", tcpHandler, clientMap)
	go suite.server.Run()
}

func (suite *tcpServerSuite) TestDial() {
	suite.T().Log("dial start")
	conn, err := net.Dial("tcp", "127.0.0.1:64202")
	if err != nil {
		suite.NoError(err, "net.Dial Error")
	}

	suite.T().Log("send authInfo")
	authInfo := &goproto.Authenticate{
		Id:    "dong5854",
		Token: "fake token",
	}

	// 인증 시작
	authInfoByte, err := proto.Marshal(authInfo)
	if err != nil {
		suite.NoError(err, "proto Marshal Error at TestHandlePacket")
	}

	_, err = conn.Write(authInfoByte)
	if err != nil {
		suite.NoError(err, "conn write Error")
	}
	suite.T().Log("sent authInfo")
	time.Sleep(500 * time.Millisecond)
	// 인증 끝

	// 메시지 수신
	go func() {
		receivedBuf := make([]byte, 1024)
		suite.T().Log("waiting to receive Message")
		n, err := conn.Read(receivedBuf)
		if err != nil {
			if err != io.EOF {
				suite.NoError(err, "conn read Error")
			}
		}

		chatParser := packets.NewChatParser()
		chatParser.Parse(receivedBuf[:n])

		suite.Equal("helloWorld", chatParser.Message())
		suite.Equal("dong5854", chatParser.SenderID())
		suite.Equal("novaeric", chatParser.ReceiverID())
		suite.Equal(handler.ChatTypeCHAT, chatParser.Type())
		suite.once.Do(func() {
			suite.finishChan <- struct{}{}
		})
	}()

	chatByte := suite.chatCreator.Create("helloWorld", "dong5854", "novaeric", FBPacket.ChatTypeCHAT)

	// 메시지 전송
	suite.T().Log("send message")
	go func() {
		for {
			suite.T().Log(chatByte)
			_, err = conn.Write(chatByte)
			if err != nil {
				suite.NoError(err, "conn write Error")
			}
			time.Sleep(1 * time.Second) // 1초마다 메시지 전송
		}
	}()

	<-suite.finishChan
}
func TestServer(t *testing.T) {
	t.Run("server test", func(t *testing.T) {
		suite.Run(t, new(tcpServerSuite))
	})
}
