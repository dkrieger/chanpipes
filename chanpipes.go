package chanpipes

import (
	"runtime"
)

// New creates a new channel and returns a read-only reference ("out") and a
// write-only reference ("in")
func New() (<-chan interface{}, chan<- interface{}) {
	out := make(chan interface{})
	in := out
	return out, in
}

// Tee takes some readable channel as input, forwarding it to a new readable
// channel ("out") immediately, yields to the scheduler, then forwarding it to
// another new readable channel ("middle") after the goroutine wakes up.
func Tee(in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	middle := make(chan interface{})
	out := make(chan interface{})
	go func(out chan<- interface{}, in <-chan interface{}, middle chan<- interface{}) {
		upstream := <-in
		out <- upstream
		runtime.Gosched()
		middle <- upstream
	}(out, in, middle)
	return out, middle
}

// FanIn is like a dynamic "select" statement for N readable channels
func FanIn(inputs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	for _, input := range inputs {
		go func(ch <-chan interface{}) {
			for {
				out <- <-ch
			}
		}(input)
	}
	return out
}
