package threadsafe

type TCPChannels struct {
	FromClient chan []byte // TODO: 클라이언트와 협의 후 데이터 타입 변경
	ToClient   chan []byte // TODO: 클라이언트와 협의 후 데이터 타입 변경
	ErrChan    chan error
}
