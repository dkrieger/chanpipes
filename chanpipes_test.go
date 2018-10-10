package chanpipes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTee(t *testing.T) {
	tees := make(map[string](<-chan interface{}))
	out, in := New()
	out, tees["foo"] = Tee(out)
	out, tees["bar"] = Tee(out)
	out, tees["baz"] = Tee(out)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(in)
	<-out
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

func TestEet(t *testing.T) {
	tees := make(map[string](chan interface{}))
	out, in := New()
	in, tees["foo"] = Eet(in)
	in, tees["bar"] = Eet(in)
	in, tees["baz"] = Eet(in)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(in)
	<-out
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

func TestEetFanIn(t *testing.T) {
	tees := []<-chan interface{}{}
	out, in := New()
	in, tee := Eet(in)
	tees = append(tees, tee)
	in, tee = Eet(in)
	tees = append(tees, tee)
	in, tee = Eet(in)
	tees = append(tees, tee)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(in)
	<-out
	output := FanIn(tees...)
	for range tees {
		assert.Equal(t, "testing", <-output)
	}
}
