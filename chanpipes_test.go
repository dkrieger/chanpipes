package chanpipes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTee(t *testing.T) {
	tees := make(map[string](chan interface{}))
	leftmost, left := New()
	right := left
	left, tees["foo"], _ = Tee(left)
	left, tees["bar"], _ = Tee(left)
	left, tees["baz"], right = Tee(left)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(right)
	<-leftmost
	for range tees {
		select {
		case msg := <-tees["foo"]:
			assert.Equal(t, "testing", msg)
		case msg := <-tees["bar"]:
			assert.Equal(t, "testing", msg)
		case msg := <-tees["baz"]:
			assert.Equal(t, "testing", msg)
		}
	}
}

func TestFanIn(t *testing.T) {
	testInput := func(input chan<- interface{}, output <-chan interface{}) {
		go func() {
			input <- "testing"
		}()
		assert.Equal(t, "testing", <-output)
	}
	foo := make(chan interface{})
	bar := make(chan interface{})
	baz := make(chan interface{})
	all := FanIn(foo, bar, baz)
	testInput(foo, all)
	testInput(bar, all)
	testInput(baz, all)
}

func TestTeeFanIn(t *testing.T) {
	tees := []<-chan interface{}{}
	leftmost, left := New()
	left, tee, _ := Tee(left)
	tees = append(tees, tee)
	left, tee, _ = Tee(left)
	tees = append(tees, tee)
	left, tee, right := Tee(left)
	tees = append(tees, tee)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(right)
	<-leftmost
	output := FanIn(tees...)
	for range tees {
		assert.Equal(t, "testing", <-output)
	}
}
