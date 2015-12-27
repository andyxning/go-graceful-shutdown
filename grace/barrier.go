package grace

import "sync/atomic"

type httpBarrier struct {
	counter int32
	Barrier chan bool
}

// Increase increases the internal counter with one
func (hb *httpBarrier) Increase() {
	atomic.AddInt32(&hb.counter, 1)
}

// Decrease decreases the internal counter with one
func (hb *httpBarrier) Decrease() {
	atomic.AddInt32(&hb.counter, -1)
}

// GetCounter returns the current internal counter
func (hb httpBarrier) GetCounter() (cur int32) {
	cur = atomic.LoadInt32(&hb.counter)
	return
}

var defaultHTTPBarrier = httpBarrier{counter: 0, Barrier: make(chan bool, 1)}
