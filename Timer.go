package help

import (
	//"log"
	"math/rand"
	"time"
)

type Timer struct {
	ti    *time.Timer
	close bool
}

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

func NewTimer() *Timer {

	return &Timer{close: false}
}

func (this *Timer) SetTimer(callback func() bool, d time.Duration, times int32) {

	go func() {

		infini := (0 == times)

		for {
			this.ti = time.NewTimer(d)
			<-this.ti.C

			if this.close {
				break
			}

			if infini {
				callback()
			} else {
				times -= 1
				callback()
				if 0 == times {
					break
				}
			}
		}

	}()
}

func (this *Timer) Killed() {

	this.close = true
	if this.ti != nil {
		this.ti.Reset(0)
	}
}
