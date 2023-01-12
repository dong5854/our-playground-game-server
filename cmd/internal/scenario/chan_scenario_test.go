package scenario

import (
	"log"
	"strconv"
	"testing"
	"time"
)

const chanBufferSize = 1000

// TODO: 설계 연습용, 어느정도 틀 잡히면 지우기
func TestChanScenario(t *testing.T) {
	ch := make(chan string, chanBufferSize)

	go dialMessage(ch)

	go func(ch chan string) {
		for i := 1; i < 10; i++ {
			go acceptMessage(i, ch)
		}
	}(ch)

	acceptMessage(1, ch)

	time.Sleep(time.Second * 10)
}

func acceptMessage(num int, ch chan string) {
	time.Sleep(time.Second * time.Duration(num))
	ch <- "message" + strconv.Itoa(num)
}

func dialMessage(ch chan string) {
	time.Sleep(time.Second * 9)
	for true {
		log.Println(<-ch)
	}
}
