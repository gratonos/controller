package controller

import (
	"sync/atomic"
)

type atomicBool struct {
	value int32
}

func (this *atomicBool) Get() bool {
	value := atomic.LoadInt32(&this.value)
	return value != 0
}

func (this *atomicBool) Set(value bool) {
	var val int32
	if value {
		val = 1
	}
	atomic.StoreInt32(&this.value, val)
}
