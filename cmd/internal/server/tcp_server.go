package server

import (
	"io"
	"log"
	"net"
)

const (
	TCP     = "tcp"
	MaxUser = 1000
)

type tcpServer struct {
	net.Listener
	// TODO: controller 에 해당하는 구조체 추가
	fromClient chan []byte // TODO: 클라이언트와 협의 후 데이터 타입 변경
	toClient   chan []byte // TODO: 클라이언트와 협의 후 데이터 타입 변경
	errChan    chan error
}

func NewTCPServer(address string) TCPServer {
	server := new(tcpServer)
	var err error
	server.Listener, err = net.Listen(TCP, address)
	if err != nil {
		log.Panic(err)
	}
	server.fromClient = make(chan []byte, MaxUser)
	server.toClient = make(chan []byte, MaxUser)
	server.errChan = make(chan error, 1)
	return server
}

func (t *tcpServer) Run() {
	defer log.Println("Stopped TCPServer")
	log.Println("Start TCPServer")
	// TODO: controller 에 해당하는 구조체가 fromClient, toClient, errChan 을 생성자 파라미터로 받아서 알맞게 처리하도록 한다. go 루틴 사용
	for {
		log.Println("waiting for TCP HandShake")
		conn, err := t.Accept()
		log.Println("successfully HandShaken")
		if err != nil {
			log.Panic(err)
		}
		go t.ReadPacket(conn)
		if len(t.errChan) != 0 {
			log.Panic(<-t.errChan)
		}
	}
}

func (t *tcpServer) Accept() (net.Conn, error) {
	conn, err := t.tcpListener.Accept() // 핸드세이크 완료까지 블로킹
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (t *tcpServer) ReadPacket(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			t.errChan <- err
		}
	}
	t.fromClient <- buf[:n]
}

func (t *tcpServer) Close() error {
	return t.tcpListener.Close()
}
