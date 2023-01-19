package handler_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	handler2 "github.com/Team-OurPlayground/our-playground-game-server/internal/handler"
	"github.com/Team-OurPlayground/our-playground-game-server/internal/util/threadsafe"
)

type tcpHandlerSuite struct {
	suite.Suite
	tcpHandler  handler2.TCPHandler
	tcpChannels *threadsafe.TCPChannels
	clientMap   *sync.Map
}

func (suite *tcpHandlerSuite) SetupSuite() {
	suite.tcpChannels = new(threadsafe.TCPChannels)
	suite.clientMap = new(sync.Map)
	suite.tcpHandler = handler2.NewTCPHandler(suite.tcpChannels, suite.clientMap)
}

func (suite *tcpHandlerSuite) TestHandlePacket() {
	go suite.tcpHandler.HandlePacket()
}

func TestAttachment(t *testing.T) {
	t.Run("tcpHandler", func(t *testing.T) {
		suite.Run(t, new(tcpHandlerSuite))
	})
}
