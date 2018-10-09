package chanpipes

import (
	"runtime"
)

func Tee(left chan interface{}) (chan interface{}, chan interface{}, chan interface{}) {
	middle := make(chan interface{})
	right := make(chan interface{})
	go func(left chan<- interface{}, right <-chan interface{}, middle chan<- interface{}) {
		upstream := <-right
		left <- upstream
		runtime.Gosched()
		middle <- upstream
	}(left, right, middle)
	return right, middle, right
}

func New() (<-chan interface{}, chan interface{}) {
	left := make(chan interface{})
	return left, left
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
