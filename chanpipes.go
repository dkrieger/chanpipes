package chanpipes

import (
	"runtime"
)

func New() (<-chan interface{}, chan<- interface{}) {
	out := make(chan interface{})
	in := out
	return out, in
}

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
