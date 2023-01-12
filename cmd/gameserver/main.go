package main

func main() {
	listener := newTCPServer("127.0.0.1:6112")
	// conn, err := listener.accept()
	defer listener.close()
}
