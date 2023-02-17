package server_test

import (
	"net"
	"sync"
	"testing"

	"github.com/Team-OurPlayground/idl/goproto"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/server"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/parser"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

type tcpServerSuite struct {
	suite.Suite
	server     server.TCPServer
	finishChan chan struct{}
	once       sync.Once
}

func (suite *tcpServerSuite) SetupSuite() {
	suite.finishChan = make(chan struct{})
	parser := parser.NewProtobufParser()
	tcpChannels := &threadsafe.TCPChannels{
		FromClient: make(chan []byte, handler.MaxUser),
		ToClient:   make(chan []byte, handler.MaxUser),
		ErrChan:    make(chan error, 1),
	}
	clientMap := new(sync.Map)

	tcpHandler := handler.NewTCPHandler(parser, tcpChannels, clientMap)
	suite.server = server.NewTCPServer("127.0.0.1:6112", tcpHandler, clientMap)
	go suite.server.Run()
}

func (suite *tcpServerSuite) TestDial() {
	suite.T().Log("dial start")
	conn, err := net.Dial("tcp", "127.0.0.1:6112")
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
	// 인증 끝

	// 메시지 수신 handler.SimulateMove
	go func() {
		receivedBuf := make([]byte, 1024)
		n, err := conn.Read(receivedBuf)
		if err != nil {
			suite.NoError(err, "conn read Error")
		}
		suite.T().Log(n)
		receivedProto := new(goproto.Data)
		if err := proto.Unmarshal(receivedBuf[:n], receivedProto); err != nil {
			suite.NoError(err, "proto unmarshal error")
		}
		suite.Equal("dong5854", receivedProto.Data)
		suite.Equal(handler.SimulateMove, receivedProto.Function)
		suite.once.Do(func() {
			suite.finishChan <- struct{}{}
		})
	}()

	echoMessage := &goproto.Data{
		Function: handler.SimulateMove,
		Data:     "dong5854",
		Dx:       1.2,
		Dy:       1.3,
	}

	echoMessageByte, err := proto.Marshal(echoMessage)
	if err != nil {
		suite.NoError(err, "proto Marshal Error")
	}

	// 메시지 전송 handler.SimulateMove
	suite.T().Log("send message")
	go func() {
		for {
			_, err = conn.Write(echoMessageByte)
			if err != nil {
				suite.NoError(err, "conn write Error")
			}
		}
	}()

	<-suite.finishChan
}
func TestServer(t *testing.T) {
	t.Run("server test", func(t *testing.T) {
		suite.Run(t, new(tcpServerSuite))
	})
}
