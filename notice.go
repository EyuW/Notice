package notice

import (
	"sync"
	"time"
)

var (
	index = 0
)

type Idempotent struct {
	mutex    sync.Mutex
	begin    time.Duration
	end      time.Duration
	dealline time.Duration

	deallineInvatil time.Duration
	delayInvatil    time.Duration
	index           int
	C               chan int
}

func NewIdempotent(delayInvatil time.Duration, deallineInvatil time.Duration) *Idempotent {
	if deallineInvatil < delayInvatil {
		deallineInvatil = delayInvatil
	}
	i := new(Idempotent)
	i.delayInvatil = delayInvatil
	i.deallineInvatil = deallineInvatil
	i.C = make(chan int)
	i.index = index
	index++

	return i
}

func (i *Idempotent) wakeup() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.begin == 0 {
		return
	}
	now := time.Duration(time.Now().Unix() * int64(time.Second))
	if now < i.end {
		return
	}
	i.begin = 0
	i.C <- i.index
}

func (i *Idempotent) Reset() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	now := time.Duration(time.Now().Unix() * int64(time.Second))
	if i.begin == 0 {
		i.begin = now
		i.end = i.begin + i.delayInvatil
		i.dealline = i.begin + i.deallineInvatil
	} else {
		if i.end >= i.dealline {
			return
		}
		end := now + i.delayInvatil
		if end >= i.dealline {
			i.end = i.dealline
		} else {
			i.end = end
		}
	}
	t := time.NewTimer(i.end - now)

	go func() {
		<-t.C
		i.wakeup()
	}()
}
