package chanpipes

import (
	"runtime"
)

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

func Eet(out chan interface{}) (chan interface{}, chan interface{}) {
	middle := make(chan interface{})
	in := make(chan interface{})
	go func(out chan<- interface{}, in <-chan interface{}, middle chan<- interface{}) {
		upstream := <-in
		out <- upstream
		runtime.Gosched()
		middle <- upstream
	}(out, in, middle)
	return in, middle
}

func New() (<-chan interface{}, chan<- interface{}) {
	out := make(chan interface{})
	in := out
	return out, in
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
