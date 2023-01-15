package threadsafe

import (
	"net"
	"sync"
)

type ClientList struct {
	sync.RWMutex
	list []net.Conn
}

func (l *ClientList) Append(c net.Conn) {
	l.Lock()
	defer l.Unlock()
	l.list = append(l.list, c)
}

func (l *ClientList) Remove(c net.Conn) {
	l.Lock()
	defer l.Unlock()
	for idx, val := range l.list {
		if val == c {
			l.list = append(l.list[:idx], l.list[idx+1:]...)
			return
		}
	}
}

func (l *ClientList) Get() []net.Conn {
	l.RLock()
	defer l.RUnlock()
	return l.list
}
