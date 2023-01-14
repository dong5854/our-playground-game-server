package threadsafe

import "sync"

type List struct {
	sync.RWMutex
	list []string
}

func (l *List) Append(i string) {
	l.Lock()
	defer l.Unlock()
	l.list = append(l.list, i)
}

func (l *List) Remove(i string) {
	l.Lock()
	defer l.Unlock()
	for idx, val := range l.list {
		if val == i {
			l.list = append(l.list[:idx], l.list[idx+1:]...)
			return
		}
	}
}

func (l *List) Get() []string {
	l.RLock()
	defer l.RUnlock()
	return l.list
}
