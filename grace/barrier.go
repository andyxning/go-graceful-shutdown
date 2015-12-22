package grace

import "sync/atomic"

type graceHTTPBarrier struct {
	counter int32
	Barrier chan bool
}

// We need to change the value of `counter` field, so the receiver of
// `Increase` and `Decrease` method will be a pointer.
func (ghb *graceHTTPBarrier) Increase() {
	atomic.AddInt32(&ghb.counter, 1)
}

func (ghb *graceHTTPBarrier) Decrease() {
	atomic.AddInt32(&ghb.counter, -1)
}

func (ghb graceHTTPBarrier) GetCounter() (cur int32) {
	cur = atomic.LoadInt32(&ghb.counter)
	return
}

var defaultGraceHTTPBarrier = graceHTTPBarrier{counter: 0, Barrier: make(chan bool, 1)}
